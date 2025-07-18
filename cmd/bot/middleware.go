package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/akionka/akionkabot/internal/service"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type ctxKey int

const (
	UserKey ctxKey = iota
	LoggerKey
)

func (b *Bot) loggerMiddleware(ctx *th.Context, update telego.Update) error {
	return withLogger(ctx, b.logger).Next(update)
}

func withLogger(ctx *th.Context, logger *TelegoLogger) *th.Context {
	return ctx.WithValue(LoggerKey, logger)
}

func getCtxLogger(ctx context.Context) (*TelegoLogger, bool) {
	logger, ok := ctx.Value(LoggerKey).(*TelegoLogger)
	return logger, ok
}

func timeElapsedMiddleware(ctx *th.Context, update telego.Update) error {
	logger, ok := getCtxLogger(ctx)
	if !ok {
		return ctx.Next(update)
	}

	startTime := time.Now()
	logger.InfoContext(ctx, "start processing update")
	ctx.UpdateID()
	err := ctx.Next(update)
	endTime := time.Now()
	logger.InfoContext(ctx, "update processed", slog.Duration("elapsed_time_ns", endTime.Sub(startTime)))
	return err
}

func userMiddleware(service *service.UserService) func(*th.Context, telego.Update) error {
	return func(ctx *th.Context, update telego.Update) error {
		fromUser := getFromUser(update)
		if fromUser == nil {
			return ctx.Next(update)
		}

		user, err := service.GetUserByTelegramID(ctx, fromUser.ID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNotFound):
				user, err = service.CreateUser(ctx, &data.User{
					TelegramID: fromUser.ID,
					Username:   fromUser.Username,
					FirstName:  fromUser.FirstName,
					LastName:   fromUser.LastName,
				})
				if err != nil {
					return err
				}
			default:
				return err
			}
		}

		return ctx.WithValue(UserKey, user).Next(update)
	}
}

func getFromUser(update telego.Update) *telego.User {
	switch {
	case update.Message != nil:
		return update.Message.From
	case update.EditedMessage != nil:
		return update.EditedMessage.From
	case update.ChannelPost != nil:
		return update.ChannelPost.From
	case update.EditedChannelPost != nil:
		return update.EditedChannelPost.From
	case update.BusinessConnection != nil:
		return &update.BusinessConnection.User
	case update.BusinessMessage != nil:
		return update.BusinessMessage.From
	case update.EditedBusinessMessage != nil:
		return update.EditedBusinessMessage.From
	case update.MessageReaction != nil:
		return update.MessageReaction.User
	case update.InlineQuery != nil:
		return &update.InlineQuery.From
	case update.ChosenInlineResult != nil:
		return &update.ChosenInlineResult.From
	case update.CallbackQuery != nil:
		return &update.CallbackQuery.From
	case update.ShippingQuery != nil:
		return &update.ShippingQuery.From
	case update.PreCheckoutQuery != nil:
		return &update.PreCheckoutQuery.From
	case update.PurchasedPaidMedia != nil:
		return &update.PurchasedPaidMedia.From
	case update.PollAnswer != nil:
		return update.PollAnswer.User
	case update.MyChatMember != nil:
		return &update.MyChatMember.From
	case update.ChatMember != nil:
		return &update.ChatMember.From
	case update.ChatJoinRequest != nil:
		return &update.ChatJoinRequest.From
	default:
		return nil
	}
}

// getCtxUser retrieves the user from the context if it exists.
// It returns the user and a boolean indicating whether the user was found in the context.
func getCtxUser(ctx context.Context) (*data.User, bool) {
	user, ok := ctx.Value(UserKey).(*data.User)
	return user, ok
}
