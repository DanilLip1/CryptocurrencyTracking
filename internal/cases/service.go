package cases

import (
	"context"
	"cryptocurrency/internal/entity"
	"strings"

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

	for i := range titles {
		titles[i] = strings.ToLower(strings.TrimSpace(titles[i]))
	}

	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "service get latest prices: add coins")
	}
	coins, err := s.repo.Get(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(err, "service get latest prices")
	}
	if len(coins) == 0 {
		return nil, errors.Wrap(entity.ErrInvalidParams, "service get latest prices: no coins")
	}
	return coins, nil
}

func (s *CoinService) GetMinPrices(ctx context.Context, titles []string) ([]entity.Coin, error) {
	for i := range titles {
		titles[i] = strings.ToLower(strings.TrimSpace(titles[i]))
	}

	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "service get min price: add coins")
	}
	coins, err := s.repo.Get(ctx, titles, WithMin())
	if err != nil {
		return nil, errors.Wrap(err, "service get min price:")
	}
	if len(coins) == 0 {
		return nil, errors.Wrap(entity.ErrInvalidParams, "service get min price: no coins")
	}
	return coins, nil
}
func (s *CoinService) GetMaxPrices(ctx context.Context, titles []string) ([]entity.Coin, error) {
	for i := range titles {
		titles[i] = strings.ToLower(strings.TrimSpace(titles[i]))
	}

	if err := s.AddCoins(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "service get max prices: add coins")
	}
	coins, err := s.repo.Get(ctx, titles, WithMax())
	if err != nil {
		return nil, errors.Wrap(err, "service get max prices")
	}
	if len(coins) == 0 {
		return nil, errors.Wrap(entity.ErrInvalidParams, "service get max prices: no coins")
	}
	return coins, nil
}

func (s *CoinService) GetPriceChangePercent(ctx context.Context, titles []string) ([]entity.Coin, error) {
	for i := range titles {
		titles[i] = strings.ToLower(strings.TrimSpace(titles[i]))
	}

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
	for i := range titles {
		titles[i] = strings.ToLower(strings.TrimSpace(titles[i]))
	}

	coins, err := s.repo.GetTitles(ctx)
	if err != nil {
		return errors.Wrap(err, "service add coins: get titles")
	}

	//missingCoins := make([]string, 0)

	for _, title := range titles {
		match := false
		for _, coin := range coins {
			if coin == title {
				match = true
				break
			}
		}
		if match {
			//missingCoins = append(missingCoins, title)
			continue
		}
		//if len(missingCoins) == 0 {
		//	return nil
		//}
		rates, err := s.provider.GetRates(ctx, []string{title})
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				continue
			}
			return errors.Wrap(err, "service add coins: get rates")
		}
		//foundCoins := make([]string, 0, len(rates))
		//for _, coin := range rates {
		//	foundCoins = append(foundCoins, coin.Title)
		//}
		if err := s.repo.AddTrackedTitles(ctx, titles); err != nil {
			return errors.Wrap(err, "service add coins: add tracked titles")
		}
		if err := s.repo.SaveCoinPrices(ctx, rates); err != nil {
			return errors.Wrap(err, "service add coins: save coin prices")
		}
		coins = append(coins, title)

	}
	return nil
}
