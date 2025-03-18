package service

import (
	"context"

	"github.com/akionka/akionkabot/internal/d2pt"
	"github.com/akionka/akionkabot/internal/data"
	"github.com/akionka/akionkabot/internal/postgres"
	"github.com/akionka/akionkabot/internal/s3"
	"github.com/akionka/akionkabot/internal/stratz"
)

type HeroProvider interface {
	GetHero(ctx context.Context, id data.HeroID) (data.Hero, error)
	// Must return heroes in the same order
	GetHeroes(ctx context.Context, ids []data.HeroID) ([]data.Hero, error)
}

var _ HeroProvider = (*postgres.HeroRepository)(nil)

type ItemProvider interface {
	GetItem(ctx context.Context, id data.ItemID) (data.Item, error)
	// Must return items in the same order
	GetItems(ctx context.Context, ids []data.ItemID) ([]data.Item, error)
}

var _ ItemProvider = (*postgres.ItemRepository)(nil)

type QuestionRepository interface {
	GetQuestion(ctx context.Context, id data.QuestionID) (*data.Question, error)
	GetQuestionAvailableForUser(ctx context.Context, userID data.UserID, isWon bool) (*data.Question, error)
	GetUserAnswer(ctx context.Context, id data.QuestionID, userID data.UserID) (*data.UserAnswer, error)
	SaveQuestion(ctx context.Context, question *data.Question) (data.QuestionID, error)
	AnswerQuestion(ctx context.Context, id data.QuestionID, userID data.UserID, answer data.UserAnswer) (data.UserAnswerID, error)
	UpdateQuestionImage(ctx context.Context, id data.QuestionID, fileID string) error
	UpdateOptionImage(ctx context.Context, id data.QuestionID, option data.Option, fileID string) error
	GetQuestionStats(ctx context.Context, id data.QuestionID) (map[data.HeroID]int, error)
}

var _ QuestionRepository = (*postgres.QuestionRepository)(nil)

type MatchRepository interface {
	GetMatch(ctx context.Context, id data.MatchID) (*data.Match, error)
	SaveMatch(ctx context.Context, match *data.Match) (data.MatchID, error)
}

var _ MatchRepository = (*postgres.MatchRepository)(nil)

type UserRepository interface {
	GetUser(ctx context.Context, id data.UserID) (*data.User, error)
	GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error)
	CreateUser(ctx context.Context, user *data.User) (data.UserID, error)
	UpdateByID(ctx context.Context, id data.UserID, updateFn func(user *data.User) (bool, error)) error
}

var _ UserRepository = (*postgres.UserRepository)(nil)

type SteamAccountProvider interface {
	GetSteamAccount(ctx context.Context, id data.SteamID) (*data.SteamAccount, error)
	GetMatchSteamAccounts(ctx context.Context, matchID data.MatchID) ([]data.SteamAccount, error)
}

type SteamAccountRepository interface {
	SaveAccount(ctx context.Context, account *data.SteamAccount) error
}

type QuestionProvider interface {
	FetchQuestion(ctx context.Context, isWon bool) (*d2pt.QuestionResponse, error)
}

var _ QuestionProvider = (*d2pt.Client)(nil)

type MatchProvider interface {
	GetMatch(ctx context.Context, id data.MatchID) (*data.Match, error)
}

var _ MatchProvider = (*stratz.Client)(nil)

type ImageFetcher interface {
	FetchImage(ctx context.Context, shortName string) (data.Image, error)
}

var _ ImageFetcher = (*s3.HeroImageFetcher)(nil)
var _ ImageFetcher = (*s3.ItemImageFetcher)(nil)
