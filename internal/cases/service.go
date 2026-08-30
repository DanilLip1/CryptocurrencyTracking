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
		return errors.Wrap(err, "service update prices: get titles")
	}
	rates, err := s.provider.GetRates(ctx, title)
	if err != nil {
		return errors.Wrap(err, "get rates")
	}
	if err = s.repo.SaveCoinPrices(ctx, rates); err != nil {
		return errors.Wrap(err, "Service update prices: save coin prices")
	}
	return nil
}

func (s *CoinService) GetLatestPrices(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entity.ErrInvalidParams, "service get latest prices: titles is empty ")
	}
	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "service get latest prices: add coins")
	}
	coins, err := s.repo.Get(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(err, "service get latest prices")
	}
	return coins, nil
}

func (s *CoinService) GetMinPrices(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "service get min price: add coins")
	}
	coins, err := s.repo.Get(ctx, titles, WithMin())
	if err != nil {
		return nil, errors.Wrap(err, "service get min price:")
	}
	return coins, nil
}
func (s *CoinService) GetMaxPrices(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "service get max prices: add coins")
	}
	coins, err := s.repo.Get(ctx, titles, WithMax())
	if err != nil {
		return nil, errors.Wrap(err, "service get max prices")
	}
	return coins, nil
}

func (s *CoinService) GetPriceChangePercent(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "service get price change percent: add coins")
	}
	coins, err := s.repo.Get(ctx, titles, WithPercent())
	if err != nil {
		return nil, errors.Wrap(err, "service get price change percent")
	}
	if len(coins) == 0 {
		return nil, errors.Wrap(entity.ErrNotFound, "service get price change percent: coins not found")
	}
	return coins, nil
}

// AddCoins добавление валюты
func (s *CoinService) AddCoins(ctx context.Context, titles []string) error {
	coins, err := s.repo.GetTitles(ctx)
	if err != nil {
		return errors.Wrap(err, "service add coins: get titles")
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
				return errors.Wrap(err, "service add coins: get rates")
			}
			if err := s.repo.AddTrackedTitles(ctx, titles); err != nil {
				return errors.Wrap(err, "service add coins: add tracked titles")
			}
			if err := s.repo.SaveCoinPrices(ctx, rates); err != nil {
				return errors.Wrap(err, "service add coins: save coin prices")
			}
		}
	}
	return nil
}
