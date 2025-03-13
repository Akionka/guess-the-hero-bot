package cache

import (
	"context"
	"fmt"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/akionka/akionkabot/internal/postgres"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type UserRepository struct {
	repo  *postgres.UserRepository
	cache *cache.Cache
}

func NewUserRepository(repo *postgres.UserRepository, cache *cache.Cache) *UserRepository {
	return &UserRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (*data.User, error) {
	userKey := fmt.Sprintf("user_%s", id)
	user := &data.User{}

	v, found := r.cache.Get(userKey)
	if found {
		user = v.(*data.User)
		return user, nil
	}

	user, err := r.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	r.cache.Set(userKey, user, cache.DefaultExpiration)

	return user, nil
}

func (r *UserRepository) GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error) {
	userKey := fmt.Sprintf("user_%d", id)
	user := &data.User{}

	v, found := r.cache.Get(userKey)
	if found {
		user = v.(*data.User)
		return user, nil
	}

	user, err := r.repo.GetUserByTelegramID(ctx, id)
	if err != nil {
		return nil, err
	}

	r.cache.Set(userKey, user, cache.DefaultExpiration)

	return user, nil

}

func (r *UserRepository) SaveUser(ctx context.Context, user *data.User) (uuid.UUID, error) {
	userKeyID := fmt.Sprintf("user_%s", user.ID)
	userKeyTelegramID := fmt.Sprintf("user_%d", user.TelegramID)

	r.cache.Delete(userKeyID)
	r.cache.Delete(userKeyTelegramID)

	userID, err := r.repo.SaveUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
