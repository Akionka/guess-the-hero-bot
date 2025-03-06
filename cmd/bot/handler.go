package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/akionka/akionkabot/data"
	"github.com/google/uuid"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoapi"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type QuestionImage struct {
	data.Image
	q   *data.Question
	buf *bytes.Reader
}

var _ telegoapi.NamedReader = (*QuestionImage)(nil)

func (qi *QuestionImage) Name() string {
	return qi.q.ID.String()
}

func (qi *QuestionImage) Read(b []byte) (int, error) {
	if qi.buf == nil {
		var buf bytes.Buffer

		if err := jpeg.Encode(&buf, qi.Image, &jpeg.Options{Quality: 100}); err != nil {
			return 0, err
		}
		qi.buf = bytes.NewReader(buf.Bytes())
	}

	return qi.buf.Read(b)
}

func (b *Bot) handleQuestionRequest(ctx *th.Context, update telego.Update) error {
	user, ok := getCtxUser(ctx)
	if !ok {
		return errors.New("no user in context")
	}

	question, err := b.questionService.GetQuestionForUser(ctx, user.ID, rand.Uint()%2 == 0)
	if err != nil {
		return err
	}

	var chatID telego.ChatID
	switch {
	case update.Message != nil:
		chatID = update.Message.Chat.ChatID()
	case update.CallbackQuery != nil:
		chatID = tu.ID(update.CallbackQuery.From.ID)
		b.AnswerCallbackQuery(ctx, tu.CallbackQuery(update.CallbackQuery.ID))
	}

	sentMsg, err := b.sendQuestion(ctx, question, chatID)
	if err != nil {
		return err
	}

	if len(question.TelegramFileID) == 0 {
		fileID := ""
		maxWidth := -1
		for _, photo := range sentMsg.Photo {
			if photo.Width > maxWidth {
				maxWidth = photo.Width
				fileID = photo.FileID
			}
		}
		return b.questionService.UpdateQuestionImage(ctx, question, fileID)
	}
	return nil

}

func (b *Bot) handleQuestionAnswer(ctx *th.Context, query telego.CallbackQuery) error {
	parts := strings.SplitN(query.Data, "_", 3)
	if len(parts) != 3 {
		return errors.New("invalid callback query")
	}

	idStr, answerStr := parts[1], parts[2]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	answer, err := strconv.Atoi(answerStr)
	if err != nil {
		return err
	}

	user, ok := getCtxUser(ctx)
	if !ok {
		return errors.New("no user in context")
	}

	question, err := b.questionService.GetQuestion(ctx, id)
	if err != nil {
		return err
	}

	var userOption *data.UserOption
	for _, o := range question.Options {
		if o.Hero.ID == answer {
			userOption = &data.UserOption{
				Option: o,
			}
			break
		}
	}

	if err = b.questionService.AnswerQuestion(ctx, user, question, userOption); err != nil {
		if errors.Is(err, data.ErrAlreadyExists) {
			return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("❌ Уже ответил на этот вопрос").WithShowAlert())
		}
		return err
	}

	isInline := query.InlineMessageID != ""
	isPrivate := query.Message != nil && query.Message.GetChat().Type == telego.ChatTypePrivate

	if isInline || !isPrivate {
		if userOption.IsCorrect {
			return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("🎉 Правильно").WithShowAlert())
		}
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Неправильно 🥀").WithShowAlert())
	}

	msg, err := b.prepareQuestionText(ctx, question, userOption)
	if err != nil {
		return err
	}

	file, err := b.prepareOptionImageFile(ctx, question, userOption)
	if err != nil {
		return err
	}

	var (
		messageID int
		chatID    telego.ChatID
	)

	if !isInline {
		messageID = query.Message.GetMessageID()
		chatID = tu.ID(query.Message.GetChat().ID)
	}

	editedMsg, err := ctx.Bot().EditMessageMedia(ctx, &telego.EditMessageMediaParams{
		ChatID:          chatID,
		MessageID:       messageID,
		InlineMessageID: query.InlineMessageID,
		Media:           tu.MediaPhoto(file),
	})
	if err != nil {
		return err
	}

	_, err = ctx.Bot().EditMessageCaption(ctx, &telego.EditMessageCaptionParams{
		ChatID:          chatID,
		MessageID:       messageID,
		InlineMessageID: query.InlineMessageID,
		Caption:         msg,
		ParseMode:       telego.ModeHTML,
		ReplyMarkup: tu.InlineKeyboard(
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("Следующий вопрос").WithCallbackData("next_question")),
			tu.InlineKeyboardRow(b.shareQuestionButton(ctx, question)),
		),
	})
	if err != nil {
		return err
	}

	if len(userOption.TelegramFileID) == 0 {
		fileID := ""
		maxWidth := -1
		for _, photo := range editedMsg.Photo {
			if photo.Width > maxWidth {
				maxWidth = photo.Width
				fileID = photo.FileID
			}
		}
		return b.questionService.UpdateOptionImage(ctx, question, &userOption.Option, fileID)
	}

	return err
}

func (b *Bot) sendQuestion(ctx *th.Context, question *data.Question, chatID telego.ChatID) (*telego.Message, error) {
	text, err := b.prepareQuestionText(ctx, question, nil)
	if err != nil {
		return nil, err
	}

	var file telego.InputFile

	if len(question.TelegramFileID) > 0 {
		file = tu.FileFromID(question.TelegramFileID)
	} else {
		file, err = b.prepareQuestionImageFile(ctx, question, nil)
		if err != nil {
			return nil, err
		}
	}

	keyboard := b.prepareQuestionKeyboard(ctx, question)

	return ctx.Bot().SendPhoto(
		ctx,
		tu.Photo(chatID, file).
			WithCaption(text).
			WithReplyMarkup(keyboard).
			WithParseMode(telego.ModeHTML),
	)
}

func (b *Bot) handleQuestionShare(ctx *th.Context, query telego.InlineQuery) error {
	qID, err := uuid.Parse(strings.TrimPrefix(strings.TrimSpace(query.Query), "question "))
	if err != nil {
		return nil
	}

	question, err := b.questionService.GetQuestion(ctx, qID)
	if err != nil {
		return err
	}

	text, err := b.prepareQuestionText(ctx, question, nil)
	if err != nil {
		return err
	}

	var file telego.InputFile

	if len(question.TelegramFileID) > 0 {
		file = tu.FileFromID(question.TelegramFileID)
	} else {
		file, err = b.prepareQuestionImageFile(ctx, question, nil)
		if err != nil {
			return err
		}
	}

	keyboard := b.prepareQuestionKeyboard(ctx, question)

	err = b.AnswerInlineQuery(
		ctx,
		tu.InlineQuery(
			query.ID,
			tu.ResultCachedPhoto(
				qID.String(),
				file.FileID,
			).
				WithCaption(text).
				WithParseMode(telego.ModeHTML).
				WithReplyMarkup(keyboard).WithTitle("Поделитьcя вопросом "+question.ID.String()),
		),
	)

	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) prepareQuestionText(ctx *th.Context, question *data.Question, userOption *data.UserOption) (string, error) {
	var correctOption data.Option
	for _, o := range question.Options {
		if o.IsCorrect {
			correctOption = o
			break
		}
	}

	var buf bytes.Buffer
	questionTempl(question.PlayerMMR, question.Items).Render(ctx, &buf)
	buf.WriteRune('\n')
	if userOption != nil {
		answerTempl(userOption.Hero, correctOption.Hero, question.PlayerPos, question.IsWon).Render(ctx, &buf)
		buf.WriteRune('\n')

		proName := ""
		if question.PlayerIsPro {
			proName = question.PlayerName
		}
		matchCredentials(correctOption.Hero.DisplayName, question.MatchID, question.PlayerID, proName).Render(ctx, &buf)
	}

	return buf.String(), nil
}

func (b *Bot) prepareQuestionImageFile(_ *th.Context, question *data.Question, userOption *data.UserOption) (telego.InputFile, error) {
	var imageFile telego.InputFile

	if len(question.TelegramFileID) > 0 {
		imageFile = tu.FileFromID(question.TelegramFileID)
	} else {
		collage, err := b.collager.Collage(question.Options, question.Items, userOption)
		if err != nil {
			return imageFile, err
		}
		imageFile = tu.File(&QuestionImage{q: question, Image: data.Image{Image: collage}})
	}

	return imageFile, nil
}

func (b *Bot) prepareOptionImageFile(_ *th.Context, question *data.Question, userOption *data.UserOption) (telego.InputFile, error) {
	var imageFile telego.InputFile

	if len(userOption.TelegramFileID) > 0 {
		imageFile = tu.FileFromID(userOption.TelegramFileID)
	} else {
		collage, err := b.collager.Collage(question.Options, question.Items, userOption)
		if err != nil {
			return imageFile, err
		}
		imageFile = tu.File(&QuestionImage{q: question, Image: data.Image{Image: collage}})
	}

	return imageFile, nil
}

func (b *Bot) prepareQuestionKeyboard(ctx *th.Context, question *data.Question) *telego.InlineKeyboardMarkup {
	btns := make([]telego.InlineKeyboardButton, 0, len(question.Options)+1)
	for _, o := range question.Options {
		cbData := fmt.Sprintf("answer_%s_%d", question.ID, o.Hero.ID)
		btns = append(btns, tu.InlineKeyboardButton(o.Hero.DisplayName).WithCallbackData(cbData))
	}
	btns = append(btns, b.shareQuestionButton(ctx, question))
	return tu.InlineKeyboard(tu.InlineKeyboardCols(2, btns...)...)
}

func (b *Bot) shareQuestionButton(_ *th.Context, question *data.Question) telego.InlineKeyboardButton {
	return tu.InlineKeyboardButton("Поделиться вопросом").WithSwitchInlineQuery("question " + question.ID.String())
}
