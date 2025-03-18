package service

import (
	"context"
	"fmt"

	"github.com/akionka/akionkabot/internal/data"
)

type MatchService struct {
	repo             MatchRepository
	itemImageFetcher ImageFetcher
}

func NewMatchService(repo MatchRepository, itemImageFetcher ImageFetcher) *MatchService {
	return &MatchService{
		repo:             repo,
		itemImageFetcher: itemImageFetcher,
	}
}

func (s *MatchService) GetMatch(ctx context.Context, id data.MatchID) (*data.Match, error) {
	match, err := s.repo.GetMatch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	return s.fetchImages(ctx, match)
}

// fetchImages fetches images for the items in the match and updates the question with the fetched images.
func (s *MatchService) fetchImages(ctx context.Context, match *data.Match) (*data.Match, error) {
	for pIdx, player := range match.Players {
		for iIdx, item := range player.Items {
			image, err := s.itemImageFetcher.FetchImage(ctx, item.ShortName)
			if err != nil {
				return nil, err
			}
			match.Players[pIdx].Items[iIdx].Image = image
		}
	}
	return match, nil
}
