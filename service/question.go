package service

import (
	"context"
	"errors"
	"time"

	"github.com/akionka/akionkabot/d2pt"
	"github.com/akionka/akionkabot/data"
	"github.com/google/uuid"
)

type QuestionService struct {
	repo             QuestionRepository
	questionFetcher  QuestionFetcher
	heroRepo         HeroRepository
	itemRepo         ItemRepository
	heroImageFetcher ImageFetcher
	itemImageFetcher ImageFetcher
}

func NewQuestionService(repo QuestionRepository, questionFetcher QuestionFetcher, heroRepo HeroRepository, itemRepo ItemRepository, heroImageFetcher ImageFetcher, itemImageFetcher ImageFetcher) *QuestionService {
	return &QuestionService{
		repo:             repo,
		questionFetcher:  questionFetcher,
		heroRepo:         heroRepo,
		itemRepo:         itemRepo,
		heroImageFetcher: heroImageFetcher,
		itemImageFetcher: itemImageFetcher,
	}
}

func (s *QuestionService) GetQuestion(ctx context.Context, id uuid.UUID) (*data.Question, error) {
	q, err := s.repo.GetQuestion(ctx, id)
	if err != nil {
		return q, nil
	}

	return s.fetchQuestionImages(ctx, q)
}

func (s *QuestionService) GetQuestionForUser(ctx context.Context, userID uuid.UUID, isWon bool) (*data.Question, error) {
	q, err := s.repo.GetQuestionAvailableForUser(ctx, userID)
	if err == nil {
		return s.fetchQuestionImages(ctx, q)
	}

	for range 4 {
		qr, err := s.questionFetcher.FetchQuestion(ctx, isWon)
		if err != nil {
			return nil, err
		}

		q = s.convertQuestionResponse(qr)

		q, err = s.repo.SaveQuestion(ctx, q)
		if err != nil {
			if errors.Is(err, data.ErrAlreadyExists) {
				continue
			}
			return nil, err
		}

		return s.fetchQuestionImages(ctx, q)
	}

	return nil, data.ErrAlreadyExists
}

func (s *QuestionService) AnswerQuestion(ctx context.Context, user *data.User, question *data.Question, userOption *data.UserOption) error {
	userOption.AnsweredAt = time.Now()
	userOption.ID = uuid.Must(uuid.NewV7())
	return s.repo.AnswerQuestion(ctx, user, question, userOption)
}

func (s *QuestionService) UpdateQuestionImage(ctx context.Context, q *data.Question, fileID string) error {
	return s.repo.UpdateQuestionImage(ctx, q, fileID)
}

func (s *QuestionService) UpdateOptionImage(ctx context.Context, q *data.Question, o *data.Option, fileID string) error {
	return s.repo.UpdateOptionImage(ctx, q, o, fileID)
}

func (s *QuestionService) fetchQuestionImages(ctx context.Context, q *data.Question) (*data.Question, error) {
	for i, option := range q.Options {
		image, err := s.heroImageFetcher.FetchImage(ctx, option.Hero.ShortName)
		if err != nil {
			return nil, err
		}
		q.Options[i].Hero.Image = image
	}

	for i, item := range q.Items {
		image, err := s.itemImageFetcher.FetchImage(ctx, item.ShortName)
		if err != nil {
			return nil, err
		}
		q.Items[i].Image = image
	}
	return q, nil
}

func (s *QuestionService) convertQuestionResponse(qr *d2pt.Question) *data.Question {
	items := []data.Item{}
	for _, item := range qr.Items {
		items = append(items, data.Item{ID: item})
	}

	options := []data.Option{}
	for _, option := range qr.WrongOptions {
		options = append(options, data.Option{
			IsCorrect: false,
			Hero:      data.Hero{ID: option},
		})
	}
	options = append(options, data.Option{
		IsCorrect: true,
		Hero:      data.Hero{ID: qr.HeroID},
	})

	return &data.Question{
		ID:          uuid.Must(uuid.NewV7()),
		MatchID:     qr.MatchID,
		PlayerID:    qr.AccountID,
		PlayerName:  qr.Name,
		PlayerIsPro: bool(qr.IsPro),
		PlayerPos:   positionFromQuestionResponse(qr.Pos),
		PlayerMMR:   qr.MMR,
		IsWon:       bool(qr.IsWon),
		CreatedAt:   time.Now(),
		Items:       items,
		Options:     options,
	}
}

func positionFromQuestionResponse(str string) data.Position {
	switch str {
	case "pos 1":
		return data.PositionCarry
	case "pos 2":
		return data.PositionMid
	case "pos 3":
		return data.PositionOfflane
	case "pos 4":
		return data.PositionSoftSupport
	case "pos 5":
		return data.PositionHardSupport
	default:
		return data.PositionUnknown
	}
}
