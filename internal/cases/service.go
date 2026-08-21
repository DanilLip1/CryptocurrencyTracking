package cases

import (
	"context"
	"cryptocurrency/internal/entity"

	"github.com/pkg/errors"
)

type CoinService struct {
	repo     Repository
	provider Provider
}

func NewCurrencyService(
	repo Repository,
	provider Provider) (*CoinService, error) {
	if repo == nil {
		return nil, errors.Wrap(entity.ErrInvalidParams, "repository is nil")
	}
	if provider == nil {
		return nil, errors.Wrap(entity.ErrInvalidParams, "provider is nil")
	}
	return &CoinService{repo: repo, provider: provider}, nil
}

func (s *CoinService) UpdatePrices(ctx context.Context) error {
	title, err := s.repo.GetTitles(ctx)
	if err != nil {
		return errors.Wrap(err, "get titles")
	}
	rates, err := s.provider.GetRates(ctx, title)
	if err != nil {
		return errors.Wrap(err, "get rates")
	}
	if err = s.repo.SaveCoinPrices(ctx, rates); err != nil {
		return errors.Wrap(err, "save coin prices")
	}
	return nil
}

func (s *CoinService) GetLatestPrices(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entity.ErrInvalidParams, "service get latest prices: titles is empty ")
	}
	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "add coins")
	}
	coins, err := s.repo.GetLatestPrices(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(err, "get latest prices")
	}
	return coins, nil
}

func (s *CoinService) GetMinPrices(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "add coins")
	}
	coins, err := s.repo.GetMinPrices(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(err, "get min prices")
	}
	return coins, nil
}
func (s *CoinService) GetMaxPrices(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "add coins")
	}
	coins, err := s.repo.GetMaxPrices(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(err, "get max prices")
	}
	return coins, nil
}

func (s *CoinService) GetPriceChangePercent(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "add coins")
	}
	coins, err := s.repo.GetPriceChangePercent(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(err, "get price change percent")
	}
	if len(coins) == 0 {
		return nil, errors.Wrap(entity.ErrNotFound, "coins not found")
	}
	return coins, nil
}

// AddCoins добавление валюты
func (s *CoinService) AddCoins(ctx context.Context, titles []string) error {
	coins, err := s.repo.GetTitles(ctx)
	if err != nil {
		return errors.Wrap(err, "get titles")
	}
	for _, title := range titles {
		match := false
		for _, coin := range coins {
			if coin == title {
				match = true
				break
			}
		}
		if !match {
			rates, err := s.provider.GetRates(ctx, titles)
			if err != nil {
				return errors.Wrap(err, "get rates")
			}
			if err := s.repo.AddTrackedTitles(ctx, titles); err != nil {
				return errors.Wrap(err, "add tracked titles")
			}
			if err := s.repo.SaveCoinPrices(ctx, rates); err != nil {
				return errors.Wrap(err, "save coin prices")
			}
		}
	}
	return nil
}
