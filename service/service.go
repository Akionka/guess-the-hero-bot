package service

import (
	"context"

	"github.com/akionka/akionkabot/s3"

	"github.com/akionka/akionkabot/d2pt"
	"github.com/akionka/akionkabot/data"
	"github.com/akionka/akionkabot/postgres"
	"github.com/google/uuid"
)

type HeroRepository interface {
	GetHeroByID(ctx context.Context, id int) (data.Hero, error)
	// Must return heroes in the same order
	GetHeroesByIDs(ctx context.Context, ids []int) ([]data.Hero, error)
}

var _ HeroRepository = (*postgres.HeroRepository)(nil)

type ItemRepository interface {
	GetItemByID(ctx context.Context, id int) (data.Item, error)
	// Must return items in the same order
	GetItemsByIDs(ctx context.Context, ids []int) ([]data.Item, error)
}

var _ ItemRepository = (*postgres.ItemRepository)(nil)

type QuestionRepository interface {
	GetQuestion(ctx context.Context, id uuid.UUID) (*data.Question, error)
	GetQuestionAvailableForUser(ctx context.Context, id uuid.UUID) (*data.Question, error)
	SaveQuestion(ctx context.Context, question *data.Question) (*data.Question, error)
	AnswerQuestion(ctx context.Context, user *data.User, question *data.Question, answer *data.UserOption) error
	UpdateQuestionImage(ctx context.Context, question *data.Question, fileID string) error
}

var _ QuestionRepository = (*postgres.QuestionRepository)(nil)

type UserRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*data.User, error)
	GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error)
	SaveUser(ctx context.Context, user *data.User) (*data.User, error)
}

var _ UserRepository = (*postgres.UserRepository)(nil)

type QuestionFetcher interface {
	FetchQuestion(ctx context.Context, isWon bool) (*d2pt.Question, error)
}

var _ QuestionFetcher = (*d2pt.Client)(nil)

type ImageFetcher interface {
	FetchImage(ctx context.Context, shortName string) (data.Image, error)
}

var _ ImageFetcher = (*s3.HeroImageFetcher)(nil)
var _ ImageFetcher = (*s3.ItemImageFetcher)(nil)
