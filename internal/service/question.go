package service

import (
	"context"
	"errors"
	"time"

	"github.com/akionka/akionkabot/internal/d2pt"
	"github.com/akionka/akionkabot/internal/data"

	"github.com/google/uuid"
)

type QuestionService struct {
	repo             QuestionRepository
	questionFetcher  QuestionFetcher
	matchFetcher     MatchFetcher
	heroRepo         HeroProvider
	itemRepo         ItemProvider
	heroImageFetcher ImageFetcher
	itemImageFetcher ImageFetcher
}

func NewQuestionService(repo QuestionRepository, questionFetcher QuestionFetcher, matchFetcher MatchFetcher, heroRepo HeroProvider, itemRepo ItemProvider, heroImageFetcher ImageFetcher, itemImageFetcher ImageFetcher) *QuestionService {
	return &QuestionService{
		repo:             repo,
		questionFetcher:  questionFetcher,
		matchFetcher:     matchFetcher,
		heroRepo:         heroRepo,
		itemRepo:         itemRepo,
		heroImageFetcher: heroImageFetcher,
		itemImageFetcher: itemImageFetcher,
	}
}

func (s *QuestionService) GetQuestion(ctx context.Context, id uuid.UUID) (*data.Question, error) {
	question, err := s.repo.GetQuestion(ctx, id)
	if err != nil {
		return question, nil
	}

	return s.fetchQuestionImages(ctx, question)
}

func (s *QuestionService) GetQuestionForUser(ctx context.Context, userID uuid.UUID, isWon bool) (*data.Question, error) {
	question, err := s.repo.GetQuestionAvailableForUser(ctx, userID, isWon)
	if err == nil {
		return s.fetchQuestionImages(ctx, question)
	}
	var errs error

	for range 4 {
		qr, err := s.questionFetcher.FetchQuestion(ctx, isWon)
		if err != nil {
			return nil, err
		}

		question = s.convertQuestionResponse(qr)
		question.ID = uuid.Must(uuid.NewV7())
		question.CreatedAt = time.Now()

		match, err := s.matchFetcher.GetMatchByID(ctx, question.Match.ID)
		if err != nil {
			return nil, err
		}

		match.AvgMMR = question.Match.AvgMMR
		question.Match = *match

		questionID, err := s.repo.SaveQuestion(ctx, question)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		question, err = s.repo.GetQuestion(ctx, questionID)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		return s.fetchQuestionImages(ctx, question)
	}

	return nil, errs
}

func (s *QuestionService) GetUserAnswer(ctx context.Context, id uuid.UUID, userID uuid.UUID) (data.UserAnswer, error) {
	return s.repo.GetUserAnswer(ctx, id, userID)
}

func (s *QuestionService) AnswerQuestion(ctx context.Context, user *data.User, question *data.Question, option data.Option) error {
	answer := data.UserAnswer{
		ID:         uuid.Must(uuid.NewV7()),
		Option:     option,
		AnsweredAt: time.Now(),
	}
	return s.repo.AnswerQuestion(ctx, user.ID, question, answer)
}

func (s *QuestionService) UpdateQuestionImage(ctx context.Context, id uuid.UUID, fileID string) error {
	return s.repo.UpdateQuestionImage(ctx, id, fileID)
}

func (s *QuestionService) UpdateOptionImage(ctx context.Context, id uuid.UUID, option data.Option, fileID string) error {
	return s.repo.UpdateOptionImage(ctx, id, option, fileID)
}

func (s *QuestionService) GetQuestionStats(ctx context.Context, questionID uuid.UUID) (map[int]int, error) {
	return s.repo.GetQuestionStats(ctx, questionID)
}

func (s *QuestionService) fetchQuestionImages(ctx context.Context, question *data.Question) (*data.Question, error) {
	for i, option := range question.Options {
		image, err := s.heroImageFetcher.FetchImage(ctx, option.Hero.ShortName)
		if err != nil {
			return nil, err
		}
		question.Options[i].Hero.Image = image
	}

	for i, item := range question.Player.Items {
		image, err := s.itemImageFetcher.FetchImage(ctx, item.ShortName)
		if err != nil {
			return nil, err
		}
		question.Player.Items[i].Image = image
	}
	return question, nil
}

func (s *QuestionService) convertQuestionResponse(qr *d2pt.QuestionResponse) *data.Question {
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
		Match: data.Match{
			ID:     qr.MatchID,
			AvgMMR: &qr.MMR,
		},
		Player: &data.MatchPlayer{
			Player: data.Player{
				SteamID: qr.AccountID,
			},
		},
		Options: options,
	}
}
