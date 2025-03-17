package main

import (
	"context"
	"log/slog"

	"github.com/akionka/akionkabot/internal/service"

	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"
)

type Bot struct {
	*telego.Bot

	logger *TelegoLogger

	collager        Collager
	questionService *service.QuestionService
	userService     *service.UserService
	playerService   *service.PlayerService
}

func NewBot(bot *telego.Bot, logger *TelegoLogger, collager Collager, questionService *service.QuestionService, userService *service.UserService, playerService *service.PlayerService) *Bot {
	return &Bot{
		Bot:             bot,
		logger:          logger,
		collager:        collager,
		questionService: questionService,
		userService:     userService,
		playerService:   playerService,
	}
}

func (b *Bot) Start(ctx context.Context) {
	const op = "Bot.Start"
	updates, _ := b.UpdatesViaLongPolling(ctx, nil)
	bh, _ := th.NewBotHandler(b.Bot, updates)

	bh.Use(b.loggerMiddleware)
	bh.Use(timeElapsedMiddleware)
	bh.Use(userMiddleware(b.userService))

	bh.Handle(b.handleQuestionRequest, th.Or(th.CommandEqual("question"), th.CommandEqual("start"), th.CallbackDataEqual("next_question")))

	bh.HandleCallbackQuery(b.handleQuestionAnswer, th.CallbackDataPrefix("answer_"))
	bh.HandleCallbackQuery(b.handleMyAnswer, th.CallbackDataPrefix("my_answer_"))
	bh.HandleCallbackQuery(b.handleStats, th.CallbackDataPrefix("stats_"))

	bh.HandleInlineQuery(b.handleQuestionShare, th.InlineQueryPrefix("question "))

	bh.HandleMessage(b.handleCmdConnect, th.CommandEqual("connect"))
	bh.HandleCallbackQuery(b.handleQueryConnectCancel, th.CallbackDataPrefix("connect_cancel"))
	bh.HandleCallbackQuery(b.handleQueryConnect, th.CallbackDataPrefix("connect"))

	go bh.Start()

	<-ctx.Done()
	b.logger.Info("stopping bot", opAttr(op))
	bh.Stop()
}

func opAttr(op string) slog.Attr {
	return slog.String("op", op)
}
