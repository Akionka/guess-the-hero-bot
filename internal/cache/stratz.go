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

func (c StratzClient) GetMatchByID(ctx context.Context, matchID int64) (*data.Match, error) {
	key := matchKey(matchID)

	v, found := c.cache.Get(key)
	if found {
		return v.(*data.Match), nil
	}

	v, err, _ := c.g.Do(key, func() (any, error) {
		return c.client.GetMatchByID(ctx, matchID)
	})
	m := v.(*data.Match)
	if err != nil {
		return v.(*data.Match), err
	}

	c.cache.Set(key, m, cache.DefaultExpiration)

	return m, err
}

func (c StratzClient) GetPlayerByID(ctx context.Context, steamID int64) (*data.SteamAccount, error) {
	key := playerKey(steamID)

	v, found := c.cache.Get(key)
	if found {
		return v.(*data.SteamAccount), nil
	}

	v, err, _ := c.g.Do(key, func() (any, error) {
		return c.client.GetPlayerByID(ctx, steamID)
	})
	p := v.(*data.SteamAccount)
	if err != nil {
		return p, err
	}

	c.cache.Set(key, p, cache.DefaultExpiration)

	return p, err
}

func matchKey(matchID int64) string {
	return fmt.Sprintf("stratz_match_%d", matchID)
}

func playerKey(steamID int64) string {
	return fmt.Sprintf("stratz_player_%d", steamID)
}
