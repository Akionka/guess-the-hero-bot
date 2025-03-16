package service

import (
	"context"

	"github.com/akionka/akionkabot/internal/data"
)

type PlayerService struct {
	provider PlayerProvider
}

func NewPlayerService(provider PlayerProvider) *PlayerService {
	return &PlayerService{
		provider: provider,
	}
}

func (s *PlayerService) GetPlayerByID(ctx context.Context, id int64) (*data.Player, error) {
	return s.provider.GetPlayerByID(ctx, id)
}
