package cases

import (
	"context"
	"cryptocurrency/internal/entity"
	"errors"
	"time"
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
	coins, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}
	var titles []string
	for _, coin := range *coins {
		titles = append(titles, coin.Title)
	}
	rates, err := s.provider.GetRates(ctx, titles)
	if err != nil {
		return err
	}
	if err = s.repo.SaveRates(ctx, rates); err != nil {
		return err
	}
	return nil
}

// GetCurrentRates возвращает последний курс по каждому запрошенному символу
func (s *CoinService) GetCurrentRates(ctx context.Context, titles []string) (*[]entity.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.New("titles is empty")
	}
	coins, err := s.repo.GetLast(ctx, titles)
	if err != nil {
		return nil, err
	}
	if coins == nil || len(*coins) == 0 {
		return nil, errors.New("coins not found")
	}
	return coins, nil
}

// CheckCoin добавление валюты
func (s *CoinService) CheckCoin(ctx context.Context, titles []string) error {
	coins, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}
	var missing []string
	for _, title := range titles {
		match := false
		for _, coin := range *coins {
			if coin.Title == title {
				match = true
				break
			}
		}
		if !match {
			missing = append(missing, title)
		}
	}
	rates, err := s.provider.GetRates(ctx, missing)
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

func (s *CoinService) Min24h(ctx context.Context, titles []string) (float64, error) {
	min24h, err := s.repo.GetMin24h(ctx, titles)
	if err != nil {
		return 0, err
	}
	return min24h, nil
}
func (s *CoinService) Max24h(ctx context.Context, titles []string) (float64, error) {
	max24h, err := s.repo.GetMax24h(ctx, titles)
	if err != nil {
		return 0, err
	}
	return max24h, nil
}

func (s *CoinService) GetPercent(ctx context.Context, titles []string, at time.Time) (float64, error) {
	at = time.Now().Add(-time.Hour)
	coins, err := s.repo.GetPercent(ctx, titles, at)
	if err != nil {
		return 0, err
	}
	if coins == nil || len(*coins) < 2 {
		return 0, errors.New("not enough data")
	}
	oldValue := (*coins)[0].Value
	NewValue := (*coins)[1].Value
	if oldValue == 0 {
		return 0, errors.New("old value is zero")
	}
	percent := ((NewValue - oldValue) / oldValue) * 100
	return percent, err
}
