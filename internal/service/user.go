package service

import (
	"context"
	"fmt"
	"time"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/google/uuid"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error) {
	user, err := s.repo.GetUserByTelegramID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *data.User) (*data.User, error) {
	user = &data.User{
		ID:         uuid.Must(uuid.NewV7()),
		TelegramID: user.TelegramID,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		CreatedAt:  time.Now(),
	}
	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	user, err = s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (s *UserService) ConnectSteamAccount(ctx context.Context, id uuid.UUID, account *data.SteamAccount) error {
	return s.repo.UpdateByID(ctx, id, func(user *data.User) (bool, error) {
		user.SteamAcc = account
		return true, nil
	})
}
