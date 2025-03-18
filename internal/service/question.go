package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/akionka/akionkabot/internal/d2pt"
	"github.com/akionka/akionkabot/internal/data"
	"github.com/google/uuid"
)

type QuestionService struct {
	repo             QuestionRepository
	matchRepo        MatchRepository
	steamAccountRepo SteamAccountRepository
	questionProvider QuestionProvider
	matchProvider    MatchProvider
	heroRepo         HeroProvider
	itemRepo         ItemProvider
	playerProvider   SteamAccountProvider
	heroImageFetcher ImageFetcher
}

func NewQuestionService(
	repo QuestionRepository, matchRepo MatchRepository, steamAccountRepo SteamAccountRepository,
	questionProvider QuestionProvider, matchProvider MatchProvider, heroRepo HeroProvider, itemRepo ItemProvider, playerProvider SteamAccountProvider,
	heroImageFetcher ImageFetcher) *QuestionService {
	return &QuestionService{
		repo:             repo,
		matchRepo:        matchRepo,
		steamAccountRepo: steamAccountRepo,
		questionProvider: questionProvider,
		matchProvider:    matchProvider,
		heroRepo:         heroRepo,
		itemRepo:         itemRepo,
		playerProvider:   playerProvider,
		heroImageFetcher: heroImageFetcher,
	}
}

func (s *QuestionService) GetQuestion(ctx context.Context, id data.QuestionID) (*data.Question, error) {
	question, err := s.repo.GetQuestion(ctx, id)
	if err != nil {
		return question, nil
	}

	return question, nil
}

// TODO: add saving all info about match gotten from stratz api
func (s *QuestionService) GetQuestionForUser(ctx context.Context, userID data.UserID, isWon bool) (*data.Question, error) {
	question, err := s.repo.GetQuestionAvailableForUser(ctx, userID, isWon)
	if err == nil {
		return s.fetchImages(ctx, question)
	}
	var errs error

	for range 4 {
		// Fetches question from D2PT API
		qr, err := s.questionProvider.FetchQuestion(ctx, isWon)
		if err != nil {
			return nil, err
		}

		question = s.convertQuestionResponse(qr)
		question.ID = data.QuestionID(uuid.Must(uuid.NewV7()))
		question.CreatedAt = time.Now()

		// Fetches match info from Stratz API
		match, err := s.matchProvider.GetMatch(ctx, question.MatchID)
		if err != nil {
			return nil, err
		}

		for i, player := range match.Players {
			itemIDs := make([]data.ItemID, 0, 6)
			for _, item := range player.Items {
				if item.ID == 0 {
					continue
				}
				itemIDs = append(itemIDs, item.ID)
			}
			items, err := s.itemRepo.GetItems(ctx, itemIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to get items %v: %w", items, err)
			}
			match.Players[i].Items = items
		}

		accounts, err := s.playerProvider.GetMatchSteamAccounts(ctx, question.MatchID)
		if err != nil {
			return nil, err
		}

		for _, acc := range accounts {
			s.steamAccountRepo.SaveAccount(ctx, &acc)
		}

		match.AvgMMR = &qr.MMR
		question.MatchID, err = s.matchRepo.SaveMatch(ctx, match)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

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

		return s.fetchImages(ctx, question)
	}

	return nil, errs
}

func (s *QuestionService) GetUserAnswer(ctx context.Context, id data.QuestionID, userID data.UserID) (*data.UserAnswer, error) {
	return s.repo.GetUserAnswer(ctx, id, userID)
}

func (s *QuestionService) AnswerQuestion(ctx context.Context, id data.QuestionID, userID data.UserID, option data.Option) (data.UserAnswerID, error) {
	answer := data.UserAnswer{
		ID:         data.UserAnswerID(uuid.Must(uuid.NewV7())),
		Option:     option,
		AnsweredAt: time.Now(),
	}
	return s.repo.AnswerQuestion(ctx, id, userID, answer)
}

func (s *QuestionService) UpdateQuestionImage(ctx context.Context, id data.QuestionID, fileID string) error {
	return s.repo.UpdateQuestionImage(ctx, id, fileID)
}

func (s *QuestionService) UpdateOptionImage(ctx context.Context, id data.QuestionID, option data.Option, fileID string) error {
	return s.repo.UpdateOptionImage(ctx, id, option, fileID)
}

func (s *QuestionService) GetQuestionStats(ctx context.Context, id data.QuestionID) (map[data.HeroID]int, error) {
	return s.repo.GetQuestionStats(ctx, id)
}

// Fetches images for heroes in the question options and updates the question with the fetched images.
func (s *QuestionService) fetchImages(ctx context.Context, question *data.Question) (*data.Question, error) {
	for i, option := range question.Options {
		image, err := s.heroImageFetcher.FetchImage(ctx, option.Hero.ShortName)
		if err != nil {
			return nil, err
		}
		question.Options[i].Hero.Image = image
	}

	return question, nil
}

func (s *QuestionService) convertQuestionResponse(qr *d2pt.QuestionResponse) *data.Question {
	options := []data.Option{}
	for _, option := range qr.WrongOptions {
		options = append(options, data.Option{
			IsCorrect: false,
			Hero:      data.Hero{ID: data.HeroID(option)},
		})
	}
	options = append(options, data.Option{
		IsCorrect: true,
		Hero:      data.Hero{ID: data.HeroID(qr.HeroID)},
	})

	return &data.Question{
		MatchID:  data.MatchID(qr.MatchID),
		PlayerID: data.SteamID(qr.SteamID),
		Options:  options,
	}
}
