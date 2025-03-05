package main

import (
	"context"
	"time"

	"github.com/akionka/akionkabot/service"
	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"
)

type Bot struct {
	*telego.Bot

	collager        Collager
	questionService *service.QuestionService
	userService     *service.UserService
}

func NewBot(bot *telego.Bot, collager Collager, questionService *service.QuestionService, userService *service.UserService) *Bot {
	return &Bot{
		bot,
		collager,
		questionService,
		userService,
	}
}

func (b *Bot) Start(ctx context.Context) {
	updates, _ := b.UpdatesViaLongPolling(ctx, nil)
	bh, _ := th.NewBotHandler(b.Bot, updates)

	bh.Use(userMiddleware(b.userService))

	bh.HandleMessage(b.handleQuestionRequest, th.Or(th.CommandEqual("question"), th.CommandEqual("start")))

	bh.HandleCallbackQuery(b.handleQuestionAnswer, th.CallbackDataPrefix("answer_"))

	done := make(chan struct{}, 1)

	go func() {
		<-ctx.Done()

		stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*20)
		defer stopCancel()

		for len(updates) > 0 {
			select {
			case <-stopCtx.Done():
				break
			case <-time.After(time.Microsecond * 100):
			}
		}

		_ = bh.StopWithContext(stopCtx)

		done <- struct{}{}
	}()

	go bh.Start()

	<-done
}
