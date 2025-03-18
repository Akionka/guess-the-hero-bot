package cache

import (
	"context"
	"fmt"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/akionka/akionkabot/internal/stratz"
	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

type StratzClient struct {
	client *stratz.Client
	cache  *cache.Cache
	g      *singleflight.Group
}

func NewCachedStratzClient(client *stratz.Client, cache *cache.Cache) StratzClient {
	return StratzClient{
		client: client,
		cache:  cache,
		g:      new(singleflight.Group),
	}
}

func (c StratzClient) GetMatch(ctx context.Context, id data.MatchID) (*data.Match, error) {
	key := matchKey(id)

	v, found := c.cache.Get(key)
	if found {
		return v.(*data.Match), nil
	}

	v, err, _ := c.g.Do(key, func() (any, error) {
		return c.client.GetMatch(ctx, id)
	})
	m := v.(*data.Match)
	if err != nil {
		return v.(*data.Match), err
	}

	c.cache.Set(key, m, cache.DefaultExpiration)

	return m, err
}

func (c StratzClient) GetSteamAccount(ctx context.Context, id data.SteamID) (*data.SteamAccount, error) {
	key := playerKey(id)

	v, found := c.cache.Get(key)
	if found {
		return v.(*data.SteamAccount), nil
	}

	v, err, _ := c.g.Do(key, func() (any, error) {
		return c.client.GetSteamAccount(ctx, id)
	})
	p := v.(*data.SteamAccount)
	if err != nil {
		return p, err
	}

	c.cache.Set(key, p, cache.DefaultExpiration)

	return p, err
}

func (c StratzClient) GetMatchSteamAccounts(ctx context.Context, id data.MatchID) ([]data.SteamAccount, error) {
	key := matchKeySteamAccounts(id)

	v, found := c.cache.Get(key)
	if found {
		return v.([]data.SteamAccount), nil
	}

	v, err, _ := c.g.Do(key, func() (any, error) {
		return c.client.GetMatchSteamAccounts(ctx, id)
	})
	steamAccounts := v.([]data.SteamAccount)
	if err != nil {
		return steamAccounts, err
	}

	c.cache.Set(key, steamAccounts, cache.DefaultExpiration)

	return steamAccounts, err
}

func matchKey(id data.MatchID) string {
	return fmt.Sprintf("stratz_match_%d", id)
}

func matchKeySteamAccounts(id data.MatchID) string {
	return fmt.Sprintf("stratz_match_steam_accounts_%d", id)
}

func playerKey(id data.SteamID) string {
	return fmt.Sprintf("stratz_player_%d", id)
}
