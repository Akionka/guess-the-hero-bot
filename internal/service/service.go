package service

import (
	"context"

	"github.com/akionka/akionkabot/internal/d2pt"
	"github.com/akionka/akionkabot/internal/data"
	"github.com/akionka/akionkabot/internal/postgres"
	"github.com/akionka/akionkabot/internal/s3"
	"github.com/akionka/akionkabot/internal/stratz"

	"github.com/google/uuid"
)

type HeroProvider interface {
	GetHeroByID(ctx context.Context, id int) (data.Hero, error)
	// Must return heroes in the same order
	GetHeroesByIDs(ctx context.Context, ids []int) ([]data.Hero, error)
}

var _ HeroProvider = (*postgres.HeroRepository)(nil)

type ItemProvider interface {
	GetItemByID(ctx context.Context, id int) (data.Item, error)
	// Must return items in the same order
	GetItemsByIDs(ctx context.Context, ids []int) ([]data.Item, error)
}

var _ ItemProvider = (*postgres.ItemRepository)(nil)

type QuestionRepository interface {
	GetQuestion(ctx context.Context, id uuid.UUID) (*data.Question, error)
	GetQuestionAvailableForUser(ctx context.Context, userID uuid.UUID, isWon bool) (*data.Question, error)
	GetUserAnswer(ctx context.Context, id uuid.UUID, userID uuid.UUID) (data.UserAnswer, error)
	SaveQuestion(ctx context.Context, question *data.Question) (uuid.UUID, error)
	AnswerQuestion(ctx context.Context, userID uuid.UUID, question *data.Question, answer data.UserAnswer) error
	UpdateQuestionImage(ctx context.Context, id uuid.UUID, fileID string) error
	UpdateOptionImage(ctx context.Context, id uuid.UUID, option data.Option, fileID string) error
	GetQuestionStats(ctx context.Context, questionID uuid.UUID) (map[int]int, error)
}

var _ QuestionRepository = (*postgres.QuestionRepository)(nil)

type UserRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*data.User, error)
	GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error)
	CreateUser(ctx context.Context, user *data.User) (uuid.UUID, error)
	UpdateByID(ctx context.Context, id uuid.UUID, updateFn func(user *data.User) (bool, error)) error
}

var _ UserRepository = (*postgres.UserRepository)(nil)

type PlayerProvider interface {
	GetPlayerByID(ctx context.Context, id int64) (*data.SteamAccount, error)
}

type QuestionFetcher interface {
	FetchQuestion(ctx context.Context, isWon bool) (*d2pt.QuestionResponse, error)
}

var _ QuestionFetcher = (*d2pt.Client)(nil)

type MatchFetcher interface {
	GetMatchByID(ctx context.Context, id int64) (*data.Match, error)
}

var _ MatchFetcher = (*stratz.Client)(nil)

type ImageFetcher interface {
	FetchImage(ctx context.Context, shortName string) (data.Image, error)
}

var _ ImageFetcher = (*s3.HeroImageFetcher)(nil)
var _ ImageFetcher = (*s3.ItemImageFetcher)(nil)
