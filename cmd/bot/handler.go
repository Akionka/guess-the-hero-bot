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

	q, err := b.questionService.GetQuestionForUser(ctx, user.ID, rand.Uint()%2 == 0)
	if err != nil {
		return err
	}

	text, err := b.prepareQuestionText(ctx, q, nil)
	if err != nil {
		return err
	}

	var file telego.InputFile

	if len(q.TelegramFileID) > 0 {
		file = tu.FileFromID(q.TelegramFileID)
	} else {
		file, err = b.prepareQuestionImageFile(ctx, q, nil)
		if err != nil {
			return err
		}
	}

	keyboard := b.prepareQuestionKeyboard(ctx, q, nil)

	var chatID telego.ChatID
	switch {
	case update.Message != nil:
		chatID = update.Message.Chat.ChatID()
	case update.CallbackQuery != nil:
		chatID = tu.ID(update.CallbackQuery.From.ID)
		b.AnswerCallbackQuery(ctx, tu.CallbackQuery(update.CallbackQuery.ID))
	}

	sentMsg, err := ctx.Bot().SendPhoto(ctx,
		tu.Photo(
			chatID,
			file,
		).
			WithCaption(text).
			WithReplyMarkup(keyboard).
			WithParseMode(telego.ModeHTML),
	)
	if err != nil {
		return err
	}

	if len(q.TelegramFileID) == 0 {
		fileID := ""
		maxWidth := -1
		for _, photo := range sentMsg.Photo {
			if photo.Width > maxWidth {
				maxWidth = photo.Width
				fileID = photo.FileID
			}
		}
		return b.questionService.UpdateQuestionImage(ctx, q, fileID)
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

	q, err := b.questionService.GetQuestion(ctx, id)
	if err != nil {
		return err
	}

	var userOption *data.UserOption
	for _, o := range q.Options {
		if o.Hero.ID == answer {
			userOption = &data.UserOption{
				Option: o,
			}
			break
		}
	}

	err = b.questionService.AnswerQuestion(ctx, user, q, userOption)
	if err != nil {
		if errors.Is(err, data.ErrAlreadyExists) {
			ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Один раз ответил - достаточно.").WithShowAlert())
			return nil
		}
		return err
	}

	msg, err := b.prepareQuestionText(ctx, q, userOption)
	if err != nil {
		return err
	}

	file, err := b.prepareOptionImageFile(ctx, q, userOption)
	if err != nil {
		return err
	}

	var (
		messageID int
		chatID    telego.ChatID
	)

	if query.Message != nil {
		messageID = query.Message.GetMessageID()
		chatID = tu.ID(query.From.ID)
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
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("Поделиться").WithSwitchInlineQuery("question "+q.ID.String())),
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
		return b.questionService.UpdateOptionImage(ctx, q, &userOption.Option, fileID)
	}

	return err
}

func (b *Bot) handleQuestionShare(ctx *th.Context, query telego.InlineQuery) error {
	qID, err := uuid.Parse(strings.TrimPrefix(strings.TrimSpace(query.Query), "question "))
	if err != nil {
		return nil
	}

	q, err := b.questionService.GetQuestion(ctx, qID)
	if err != nil {
		return err
	}

	text, err := b.prepareQuestionText(ctx, q, nil)
	if err != nil {
		return err
	}

	var file telego.InputFile

	if len(q.TelegramFileID) > 0 {
		file = tu.FileFromID(q.TelegramFileID)
	} else {
		file, err = b.prepareQuestionImageFile(ctx, q, nil)
		if err != nil {
			return err
		}
	}

	keyboard := b.prepareQuestionKeyboard(ctx, q, nil)

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
				WithReplyMarkup(keyboard).WithTitle("Поделитьcя вопросом "+q.ID.String()),
		),
	)

	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) prepareQuestionText(ctx *th.Context, q *data.Question, userOption *data.UserOption) (string, error) {
	var correctOption data.Option
	for _, o := range q.Options {
		if o.IsCorrect {
			correctOption = o
			break
		}
	}

	var buf bytes.Buffer
	questionTempl(q.PlayerMMR, q.Items).Render(ctx, &buf)
	buf.WriteRune('\n')
	if userOption != nil {
		answerTempl(userOption.Hero, correctOption.Hero, q.PlayerPos, q.IsWon).Render(ctx, &buf)
		buf.WriteRune('\n')

		proName := ""
		if q.PlayerIsPro {
			proName = q.PlayerName
		}
		matchCredentials(correctOption.Hero.DisplayName, q.MatchID, q.PlayerID, proName).Render(ctx, &buf)
	}

	return buf.String(), nil
}

func (b *Bot) prepareQuestionImageFile(_ *th.Context, q *data.Question, userOption *data.UserOption) (telego.InputFile, error) {
	var imageFile telego.InputFile

	if len(q.TelegramFileID) > 0 {
		imageFile = tu.FileFromID(q.TelegramFileID)
	} else {
		collage, err := b.collager.Collage(q.Options, q.Items, userOption)
		if err != nil {
			return imageFile, err
		}
		imageFile = tu.File(&QuestionImage{q: q, Image: data.Image{Image: collage}})
	}

	return imageFile, nil
}

func (b *Bot) prepareOptionImageFile(_ *th.Context, q *data.Question, userOption *data.UserOption) (telego.InputFile, error) {
	var imageFile telego.InputFile

	if len(userOption.TelegramFileID) > 0 {
		imageFile = tu.FileFromID(userOption.TelegramFileID)
	} else {
		collage, err := b.collager.Collage(q.Options, q.Items, userOption)
		if err != nil {
			return imageFile, err
		}
		imageFile = tu.File(&QuestionImage{q: q, Image: data.Image{Image: collage}})
	}

	return imageFile, nil
}

func (b *Bot) prepareQuestionKeyboard(_ *th.Context, q *data.Question, userOption *data.Option) *telego.InlineKeyboardMarkup {
	btns := make([]telego.InlineKeyboardButton, len(q.Options))
	for i, o := range q.Options {
		cbData := fmt.Sprintf("answer_%s_%d", q.ID, o.Hero.ID)
		btns[i] = tu.InlineKeyboardButton(o.Hero.DisplayName).WithCallbackData(cbData)
	}
	return tu.InlineKeyboard(tu.InlineKeyboardCols(2, btns...)...)
}
