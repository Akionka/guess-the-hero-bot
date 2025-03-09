package service

import (
	"context"
	"fmt"
	"time"

	"github.com/akionka/akionkabot/internal/data"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type UserService struct {
	repo  UserRepository
	cache *cache.Cache
}

func NewUserService(repo UserRepository, cache *cache.Cache) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

func (s UserService) GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error) {
	userKey := fmt.Sprintf("user_%d", id)
	user := &data.User{}

	v, found := s.cache.Get(userKey)
	if found {
		user = v.(*data.User)
		return user, nil
	}

	user, err := s.repo.GetUserByTelegramID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.cache.Set(userKey, user, cache.DefaultExpiration)

	return user, nil
}

func (s UserService) CreateUser(ctx context.Context, user *data.User) (*data.User, error) {
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
