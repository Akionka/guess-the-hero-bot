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
	key := userKey(id)
	user := &data.User{}

	v, found := r.cache.Get(key)
	if found {
		user = v.(*data.User)
		return user, nil
	}

	user, err := r.repo.GetUser(ctx, id)
	if err != nil {
		return user, err
	}

	r.cache.Set(key, user, cache.DefaultExpiration)

	return user, nil
}

func (r *UserRepository) GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error) {
	key := userTgKey(id)
	user := &data.User{}

	v, found := r.cache.Get(key)
	if found {
		user = v.(*data.User)
		return user, nil
	}

	user, err := r.repo.GetUserByTelegramID(ctx, id)
	if err != nil {
		return user, err
	}

	r.cache.Set(key, user, cache.DefaultExpiration)

	return user, nil

}

func (r *UserRepository) CreateUser(ctx context.Context, user *data.User) (uuid.UUID, error) {
	userIDKey := userKey(user.ID)
	userTgIDKey := userTgKey(user.TelegramID)

	r.cache.Delete(userIDKey)
	r.cache.Delete(userTgIDKey)

	userID, err := r.repo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (r *UserRepository) UpdateByID(ctx context.Context, id uuid.UUID, updateFn func(user *data.User) (bool, error)) error {
	userKeyID := userKey(id)

	v, found := r.cache.Get(userKeyID)
	if found {
		user := v.(*data.User)
		userTgIDKey := userTgKey(user.TelegramID)
		r.cache.Delete(userTgIDKey)
	}

	r.cache.Delete(userKeyID)

	return r.repo.UpdateByID(ctx, id, updateFn)
}

func userKey(id uuid.UUID) string {
	return fmt.Sprintf("user_%s", id)
}

func userTgKey(id int64) string {
	return fmt.Sprintf("user_tg_%d", id)
}
