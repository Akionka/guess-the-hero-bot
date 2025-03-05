package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
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

func (b *Bot) handleQuestionRequest(ctx *th.Context, update telego.Message) error {
	user, ok := getCtxUser(ctx)
	if !ok {
		return errors.New("no user in context")
	}

	q, err := b.questionService.GetQuestionForUser(ctx, user.ID, true)
	if err != nil {
		return err
	}

	msg, err := b.prepareQuestionMessage(ctx, q, nil)
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

	sentMsg, err := ctx.Bot().SendPhoto(ctx,
		tu.Photo(
			update.Chat.ChatID(),
			file,
		).
			WithCaption(msg).
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

	_ = user
	err = b.questionService.AnswerQuestion(ctx, user, q, userOption)
	if err != nil {
		return err
	}

	msg, err := b.prepareQuestionMessage(ctx, q, userOption)
	if err != nil {
		return err
	}

	file, err := b.prepareQuestionImageFile(ctx, q, userOption)
	if err != nil {
		return err
	}

	_, err = ctx.Bot().EditMessageMedia(ctx, &telego.EditMessageMediaParams{
		ChatID:    tu.ID(query.From.ID),
		MessageID: query.Message.GetMessageID(),
		Media:     tu.MediaPhoto(file),
	})
	if err != nil {
		return err
	}
	_, err = ctx.Bot().EditMessageCaption(ctx, &telego.EditMessageCaptionParams{
		ChatID:          tu.ID(query.From.ID),
		MessageID:       query.Message.GetMessageID(),
		InlineMessageID: query.InlineMessageID,
		Caption:         msg,
		ParseMode:       telego.ModeHTML,
		// ReplyMarkup:     tu.InlineKeyboard(tu.InlineKeyboardRow(tu.InlineKeyboardButton("Следующий").WithCallbackData("question_next"))),
	})
	if err != nil {
		return err
	}

	return err
}

func (b *Bot) prepareQuestionMessage(ctx *th.Context, q *data.Question, userOption *data.UserOption) (string, error) {
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

func (b *Bot) prepareQuestionKeyboard(_ *th.Context, q *data.Question, userOption *data.Option) *telego.InlineKeyboardMarkup {
	btns := make([]telego.InlineKeyboardButton, len(q.Options))
	for i, o := range q.Options {
		cbData := fmt.Sprintf("answer_%s_%d", q.ID, o.Hero.ID)
		btns[i] = tu.InlineKeyboardButton(o.Hero.DisplayName).WithCallbackData(cbData)
	}
	return tu.InlineKeyboard(tu.InlineKeyboardCols(2, btns...)...)
}
