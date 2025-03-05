package main

import (
	"context"
	"errors"

	"github.com/akionka/akionkabot/data"
	"github.com/akionka/akionkabot/service"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type ctxKey int

const (
	UserKey ctxKey = iota
)

func userMiddleware(service *service.UserService) func(*th.Context, telego.Update) error {
	return func(ctx *th.Context, update telego.Update) error {
		var fromUser *telego.User
		if update.Message != nil {
			fromUser = update.Message.From
		}
		if update.EditedMessage != nil {
			fromUser = update.EditedMessage.From
		}
		if update.ChannelPost != nil {
			fromUser = update.ChannelPost.From
		}
		if update.EditedChannelPost != nil {
			fromUser = update.EditedChannelPost.From
		}
		if update.BusinessConnection != nil {
			fromUser = &update.BusinessConnection.User
		}
		if update.BusinessMessage != nil {
			fromUser = update.BusinessMessage.From
		}
		if update.EditedBusinessMessage != nil {
			fromUser = update.EditedBusinessMessage.From
		}
		if update.MessageReaction != nil {
			fromUser = update.MessageReaction.User
		}
		if update.InlineQuery != nil {
			fromUser = &update.InlineQuery.From
		}
		if update.ChosenInlineResult != nil {
			fromUser = &update.ChosenInlineResult.From
		}
		if update.CallbackQuery != nil {
			fromUser = &update.CallbackQuery.From
		}
		if update.ShippingQuery != nil {
			fromUser = &update.ShippingQuery.From
		}
		if update.PreCheckoutQuery != nil {
			fromUser = &update.PreCheckoutQuery.From
		}
		if update.PurchasedPaidMedia != nil {
			fromUser = &update.PurchasedPaidMedia.From
		}
		if update.PollAnswer != nil {
			fromUser = update.PollAnswer.User
		}
		if update.MyChatMember != nil {
			fromUser = &update.MyChatMember.From
		}
		if update.ChatMember != nil {
			fromUser = &update.ChatMember.From
		}
		if update.ChatJoinRequest != nil {
			fromUser = &update.ChatJoinRequest.From
		}

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

func getCtxUser(ctx context.Context) (*data.User, bool) {
	user, ok := ctx.Value(UserKey).(*data.User)
	return user, ok
}
