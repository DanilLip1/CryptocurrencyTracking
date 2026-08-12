package cases

import (
	"context"
	"cryptocurrency/internal/entity"
	"errors"
)

type CoinService struct {
	repo     Repository
	provider ExchangeRateProvider
}

func NewCurrencyService(
	repo Repository,
	provider ExchangeRateProvider) (*CoinService, error) {
	if repo == nil {
		return nil, errors.New("repository is nil")
	}
	if provider == nil {
		return nil, errors.New("exchange rate provider is nil")
	}
	return &CoinService{repo: repo, provider: provider}, nil
}

func (s *CoinService) UpdateRates(ctx context.Context) error {
	title, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}
	rates, err := s.provider.GetRates(ctx, title)
	if err != nil {
		return err
	}
	if err = s.repo.SaveRates(ctx, rates); err != nil {
		return err
	}
	return nil
}

// GetLast возвращает последний курс по каждому запрошенному символу
func (s *CoinService) GetLast(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.New("titles is empty")
	}
	if err := s.AddCoin(ctx, titles); err != nil {
		return nil, err
	}
	coins, err := s.repo.GetLast(ctx, titles)
	if err != nil {
		return nil, err
	}
	return coins, nil
}

func (s *CoinService) GetMinPrice(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if err := s.AddCoin(ctx, titles); err != nil {
		return nil, err
	}
	coins, err := s.repo.GetMinPrice(ctx, titles)
	if err != nil {
		return nil, err
	}
	return coins, nil
}
func (s *CoinService) GetMaxPrice(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if err := s.AddCoin(ctx, titles); err != nil {
		return nil, err
	}
	coins, err := s.repo.GetMaxPrice(ctx, titles)
	if err != nil {
		return nil, err
	}
	return coins, nil
}

func (s *CoinService) GetChangePercent(ctx context.Context, titles []string) ([]entity.Coin, error) {
	if err := s.AddCoin(ctx, titles); err != nil {
		return nil, err
	}
	coins, err := s.repo.GetChangePercent(ctx, titles)
	if err != nil {
		return nil, err
	}
	if len(coins) == 0 {
		return nil, errors.New("coins not found")
	}
	return coins, nil
}

// AddCoin добавление валюты
func (s *CoinService) AddCoin(ctx context.Context, titles []string) error {
	coins, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
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
			coins = append(coins, title)
		}
	}
	rates, err := s.provider.GetRates(ctx, coins)
	if err != nil {
		return err
	}

	if err := s.repo.CreateCoin(ctx, rates); err != nil {
		return err
	}
	if err := s.UpdateRates(ctx); err != nil {
		return err
	}
	return nil
}
