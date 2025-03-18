package service

import (
	"context"

	"github.com/akionka/akionkabot/internal/data"
)

type SteamAccountService struct {
	provider SteamAccountProvider
}

func NewSteamAccountService(provider SteamAccountProvider) *SteamAccountService {
	return &SteamAccountService{
		provider: provider,
	}
}

func (s *SteamAccountService) GetSteamAccount(ctx context.Context, id data.SteamID) (*data.SteamAccount, error) {
	return s.provider.GetSteamAccount(ctx, id)
}
