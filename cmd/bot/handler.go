package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/akionka/akionkabot/internal/data"

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
	match, err := b.matchService.GetMatch(ctx, question.MatchID)
	if err != nil {
		return err
	}
	steamAccount, err := b.playerService.GetSteamAccount(ctx, question.PlayerID)
	if err != nil {
		return err
	}
	var player data.Player
	for _, p := range match.Players {
		if p.SteamAccountID == question.PlayerID {
			player = p
			break
		}
	}

	var chatID telego.ChatID
	switch {
	case update.Message != nil:
		chatID = update.Message.Chat.ChatID()
	case update.CallbackQuery != nil:
		chatID = tu.ID(update.CallbackQuery.From.ID)
		b.AnswerCallbackQuery(ctx, tu.CallbackQuery(update.CallbackQuery.ID))
	}

	sentMsg, err := b.sendQuestion(ctx, question, &player, match, steamAccount, chatID)
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
		return b.questionService.UpdateQuestionImage(ctx, question.ID, fileID)
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

	question, err := b.questionService.GetQuestion(ctx, data.QuestionID(id))
	if err != nil {
		return err
	}
	match, err := b.matchService.GetMatch(ctx, question.MatchID)
	if err != nil {
		return err
	}
	steamAccount, err := b.playerService.GetSteamAccount(ctx, question.PlayerID)
	if err != nil {
		return err
	}
	var player data.Player
	for _, p := range match.Players {
		if p.SteamAccountID == question.PlayerID {
			player = p
			break
		}
	}

	var (
		userAnswer    data.Option
		correctOption data.Option
	)

	for _, opt := range question.Options {
		if opt.Hero.ID == data.HeroID(answer) {
			userAnswer = opt
		}
		if opt.IsCorrect {
			correctOption = opt
		}
	}

	if _, err = b.questionService.AnswerQuestion(ctx, question.ID, user.ID, userAnswer); err != nil {
		if errors.Is(err, data.ErrAlreadyExists) {
			return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("❌ Уже ответил на этот вопрос").WithShowAlert())
		}
		return err
	}

	isInline := query.InlineMessageID != ""
	isPrivate := query.Message != nil && query.Message.GetChat().Type == telego.ChatTypePrivate

	if isInline || !isPrivate {
		if userAnswer.IsCorrect {
			return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("🎉 Правильно").WithShowAlert())
		}
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText(fmt.Sprintf("🥀 Неправильно.\nНа самом деле это %s.", correctOption.Hero.DisplayName)).WithShowAlert())
	}

	msg, err := b.questionText(ctx, question, match, steamAccount, &userAnswer)
	if err != nil {
		return err
	}

	file, err := b.optionImageFile(ctx, question, &player, &userAnswer)
	if err != nil {
		return err
	}

	chatID := tu.ID(query.Message.GetChat().ID)
	messageID := query.Message.GetMessageID()

	editedMsg, err := ctx.Bot().EditMessageMedia(ctx, &telego.EditMessageMediaParams{
		ChatID:    chatID,
		MessageID: messageID,
		Media:     tu.MediaPhoto(file),
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

	if len(userAnswer.TelegramFileID) == 0 {
		fileID := ""
		maxWidth := -1
		for _, photo := range editedMsg.Photo {
			if photo.Width > maxWidth {
				maxWidth = photo.Width
				fileID = photo.FileID
			}
		}
		return b.questionService.UpdateOptionImage(ctx, question.ID, userAnswer, fileID)
	}

	return err
}

func (b *Bot) sendQuestion(ctx *th.Context, question *data.Question, player *data.Player, match *data.Match, steamAccount *data.SteamAccount, chatID telego.ChatID) (*telego.Message, error) {
	text, err := b.questionText(ctx, question, match, steamAccount, nil)
	if err != nil {
		return nil, err
	}

	var file telego.InputFile

	if len(question.TelegramFileID) > 0 {
		file = tu.FileFromID(question.TelegramFileID)
	} else {
		file, err = b.questionImageFile(ctx, question, player)
		if err != nil {
			return nil, err
		}
	}

	keyboard := b.questionKeyboard(ctx, question)

	return ctx.Bot().SendPhoto(
		ctx,
		tu.Photo(chatID, file).
			WithCaption(text).
			WithReplyMarkup(keyboard).
			WithParseMode(telego.ModeHTML),
	)
}

func (b *Bot) handleQuestionShare(ctx *th.Context, query telego.InlineQuery) error {
	id, err := uuid.Parse(strings.TrimPrefix(strings.TrimSpace(query.Query), "question "))
	if err != nil {
		return nil
	}
	qID := data.QuestionID(id)

	question, err := b.questionService.GetQuestion(ctx, qID)
	if err != nil {
		return err
	}
	match, err := b.matchService.GetMatch(ctx, question.MatchID)
	if err != nil {
		return err
	}
	steamAccount, err := b.playerService.GetSteamAccount(ctx, question.PlayerID)
	if err != nil {
		return err
	}
	var player data.Player
	for _, p := range match.Players {
		if p.SteamAccountID == question.PlayerID {
			player = p
			break
		}
	}

	text, err := b.questionText(ctx, question, match, steamAccount, nil)
	if err != nil {
		return err
	}

	var file telego.InputFile

	if len(question.TelegramFileID) > 0 {
		file = tu.FileFromID(question.TelegramFileID)
	} else {
		file, err = b.questionImageFile(ctx, question, &player)
		if err != nil {
			return err
		}
	}

	keyboard := b.questionKeyboard(ctx, question)

	err = b.AnswerInlineQuery(
		ctx,
		tu.InlineQuery(
			query.ID,
			tu.ResultCachedPhoto(
				id.String(),
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

func (b *Bot) handleMyAnswer(ctx *th.Context, query telego.CallbackQuery) error {
	id, err := uuid.Parse(strings.TrimPrefix(query.Data, "my_answer_"))
	if err != nil {
		return err
	}
	qID := data.QuestionID(id)

	user, ok := getCtxUser(ctx)
	if !ok {
		return errors.New("no user in context")
	}

	question, err := b.questionService.GetQuestion(ctx, qID)
	if err != nil {
		return err
	}

	answer, err := b.questionService.GetUserAnswer(ctx, qID, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("❌ Ты ещё не ответил на этот вопрос.").WithShowAlert())
		}
		return err
	}

	var correctAnswer data.Option
	for _, option := range question.Options {
		if option.IsCorrect {
			correctAnswer = option
			break
		}
	}

	var emoji string
	if answer.IsCorrect {
		emoji = "✅"
	} else {
		emoji = "❌"
	}

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText(fmt.Sprintf("Твой ответ: %s\nПравильный ответ: %s %s", answer.Hero.DisplayName, correctAnswer.Hero.DisplayName, emoji)).WithShowAlert())
}

func (b *Bot) handleStats(ctx *th.Context, query telego.CallbackQuery) error {
	user, ok := getCtxUser(ctx)
	if !ok {
		return errors.New("no user in context")
	}

	id, err := uuid.Parse(strings.TrimPrefix(query.Data, "stats_"))
	if err != nil {
		return err
	}
	qID := data.QuestionID(id)

	userAnswer, err := b.questionService.GetUserAnswer(ctx, qID, user.ID)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText("Статистика доступна только для ответивших"))
		}
		return err
	}

	question, err := b.questionService.GetQuestion(ctx, qID)
	if err != nil {
		return err
	}

	stats, err := b.questionService.GetQuestionStats(ctx, qID)
	if err != nil {
		return err
	}

	var statsText strings.Builder
	statsText.WriteString("📊 Статистика ответов\n\n")

	totalAnswers := 0
	for _, count := range stats {
		totalAnswers += count
	}

	if totalAnswers == 0 {
		statsText.WriteString("На этот вопрос ещё никто не ответил.")
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText(statsText.String()).WithShowAlert())
	}

	for _, opt := range question.Options {
		count := stats[opt.Hero.ID]
		percentage := float64(count) / float64(totalAnswers) * 100

		indicator := ""
		if opt.IsCorrect {
			indicator = "✅ "
		}
		if userAnswer.Hero.ID == opt.Hero.ID {
			indicator += "👤 "
		}

		progressBar := createProgressBar(percentage)

		statsText.WriteString(fmt.Sprintf("%s%s: %d (%.1f%%)\n%s\n",
			indicator, opt.Hero.DisplayName, count, percentage, progressBar))
	}

	statsText.WriteString(fmt.Sprintf("\nВсего ответов: %d", totalAnswers))

	fmt.Printf("statsText.Len(): %v\n", statsText.Len())

	return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID).WithText(statsText.String()).WithShowAlert())
}

func createProgressBar(percentage float64) string {
	const barWidth = 10
	numFilled := min(int((percentage/100)*float64(barWidth)), barWidth)

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < numFilled {
			bar += "■"
		} else {
			bar += "□"
		}
	}
	bar += "]"

	return bar
}

func (b *Bot) questionText(ctx *th.Context, question *data.Question, match *data.Match, steamAccount *data.SteamAccount, userAnswer *data.Option) (string, error) {
	var correctOption data.Option
	for _, opt := range question.Options {
		if opt.IsCorrect {
			correctOption = opt
			break
		}
	}

	var buf bytes.Buffer
	avgMMR := 0
	if match.AvgMMR != nil {
		avgMMR = *match.AvgMMR
	}

	var player data.Player
	for _, p := range match.Players {
		if p.SteamAccountID == question.PlayerID {
			player = p
			break
		}
	}

	questionTempl(avgMMR, player.Items).Render(ctx, &buf)
	buf.WriteRune('\n')
	if userAnswer != nil {
		answerTempl(userAnswer.Hero, correctOption.Hero, player.Position, match.PlayerWon(player.SteamAccountID)).Render(ctx, &buf)
		buf.WriteRune('\n')

		proName := steamAccount.ProName
		matchCredentials(correctOption.Hero.DisplayName, int64(question.MatchID), int64(question.PlayerID), proName).Render(ctx, &buf)
	}

	return buf.String(), nil
}

func (b *Bot) questionImageFile(_ *th.Context, question *data.Question, player *data.Player) (telego.InputFile, error) {
	var imageFile telego.InputFile

	if len(question.TelegramFileID) > 0 {
		imageFile = tu.FileFromID(question.TelegramFileID)
	} else {
		collage, err := b.collager.Collage(question.Options, player.Items, nil)
		if err != nil {
			return imageFile, err
		}
		imageFile = tu.File(&QuestionImage{q: question, Image: data.Image{Image: collage}})
	}

	return imageFile, nil
}

func (b *Bot) optionImageFile(_ *th.Context, question *data.Question, player *data.Player, userAnswer *data.Option) (telego.InputFile, error) {
	var imageFile telego.InputFile

	if len(userAnswer.TelegramFileID) > 0 {
		imageFile = tu.FileFromID(userAnswer.TelegramFileID)
	} else {
		collage, err := b.collager.Collage(question.Options, player.Items, userAnswer)
		if err != nil {
			return imageFile, err
		}
		imageFile = tu.File(&QuestionImage{q: question, Image: data.Image{Image: collage}})
	}

	return imageFile, nil
}

func (b *Bot) questionKeyboard(ctx *th.Context, question *data.Question) *telego.InlineKeyboardMarkup {
	btns := make([]telego.InlineKeyboardButton, 0, len(question.Options)+3)
	for _, opt := range question.Options {
		cb := fmt.Sprintf("answer_%s_%d", question.ID, opt.Hero.ID)
		btns = append(btns, tu.InlineKeyboardButton(opt.Hero.DisplayName).WithCallbackData(cb))
	}
	myAnswerCb := fmt.Sprintf("my_answer_%s", question.ID)
	statsCb := fmt.Sprintf("stats_%s", question.ID)
	btns = append(btns, tu.InlineKeyboardButton("Мой ответ").WithCallbackData(myAnswerCb), tu.InlineKeyboardButton("Статистика").WithCallbackData(statsCb))
	btns = append(btns, b.shareQuestionButton(ctx, question))
	return tu.InlineKeyboard(tu.InlineKeyboardCols(2, btns...)...)
}

func (b *Bot) shareQuestionButton(_ *th.Context, question *data.Question) telego.InlineKeyboardButton {
	return tu.InlineKeyboardButton("Поделиться вопросом").WithSwitchInlineQuery("question " + question.ID.String())
}

func (b *Bot) handleCmdConnect(ctx *th.Context, message telego.Message) error {
	const usage = "Использование: /connect <id>\n\nСвой ID можно узнать на dotabuff.com, stratz.com или в профиле Dota 2."
	_, _, args := tu.ParseCommand(message.Text)

	if len(args) < 1 {
		_, err := ctx.Bot().SendMessage(ctx, tu.Message(message.Chat.ChatID(), usage))
		return err
	}

	steamID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		_, err := ctx.Bot().SendMessage(ctx, tu.Message(message.Chat.ChatID(), usage))
		return err
	}
	playerID := data.SteamID(steamID)

	user, ok := getCtxUser(ctx)
	if !ok {
		return errors.New("no user in context")
	}

	acc, err := b.playerService.GetSteamAccount(ctx, playerID)
	if err != nil {
		return err
	}

	ctx.Bot().SendMessage(ctx, tu.Messagef(
		message.Chat.ChatID(),
		"Подключаем?\n\nSteamID: %d\nИмя: %s",
		acc.ID, acc.Name,
	).
		WithReplyMarkup(tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("Да").WithCallbackData(fmt.Sprintf("connect_%d_%d", user.TelegramID, acc.ID)),
				tu.InlineKeyboardButton("Нет").WithCallbackData(fmt.Sprintf("connect_cancel_%d", user.TelegramID)),
			),
		)),
	)

	return err
}

func (b *Bot) handleQueryConnectCancel(ctx *th.Context, query telego.CallbackQuery) error {
	idStr, found := strings.CutPrefix(query.Data, "connect_cancel_")
	if !found {
		return errors.New("could not cut prefix for connect cancel")
	}

	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return fmt.Errorf("failed to parse id: %w", err)
	}

	user, ok := getCtxUser(ctx)
	if !ok {
		return errors.New("no user in context")
	}

	if user.TelegramID != id {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(
			query.ID,
		).WithText("не для тебя"))
	}

	chatID := tu.ID(query.Message.GetChat().ID)
	msgID := query.Message.GetMessageID()

	_, err = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(chatID, msgID, "Подключение отменено"))
	return err
}

func (b *Bot) handleQueryConnect(ctx *th.Context, query telego.CallbackQuery) error {
	rest, found := strings.CutPrefix(query.Data, "connect_")
	if !found {
		return errors.New("could not cut prefix for connect")
	}

	strs := strings.SplitN(rest, "_", 2)
	if len(strs) != 2 {
		return fmt.Errorf("not enough arguments for connect got %d, expected %d", len(strs), 2)
	}

	idStr, steamIDStr := strs[0], strs[1]
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return fmt.Errorf("failed to parse id: %w", err)
	}
	steamID, err := strconv.ParseInt(steamIDStr, 10, 0)
	if err != nil {
		return fmt.Errorf("failed to parse steam id: %w", err)
	}
	playerID := data.SteamID(steamID)

	user, ok := getCtxUser(ctx)
	if !ok {
		return errors.New("no user in context")
	}

	if user.TelegramID != id {
		return ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(
			query.ID,
		).WithText("не для тебя"))
	}

	acc, err := b.playerService.GetSteamAccount(ctx, playerID)
	if err != nil {
		return fmt.Errorf("failed to get player account by id: %w", err)
	}

	chatID := tu.ID(query.Message.GetChat().ID)
	msgID := query.Message.GetMessageID()

	if err := b.userService.ConnectSteamAccount(ctx, user.ID, acc.ID); err != nil {
		_, err := ctx.Bot().EditMessageText(ctx, tu.EditMessageText(chatID, msgID, "Ошибка подключения.\n\nПовторим?").WithReplyMarkup(tu.InlineKeyboard(
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("Повторить попытку").WithCallbackData(fmt.Sprintf("connect_%d_%d", user.TelegramID, acc.ID))),
		)))
		return err
	}

	_, err = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(chatID, msgID, "Подключено"))

	return err
}
