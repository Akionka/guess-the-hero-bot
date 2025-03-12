package service

import (
	"context"
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
	user, err := s.repo.SaveUser(ctx, user)
	return user, err
}
