package service

import (
	"context"
	"cryptocurrency/internal/dto"
	"cryptocurrency/internal/entity"
	"errors"
	"time"
)

type CurrencyService struct {
	cryptoRepo entity.CryptoCurrencyRepository
	rateRepo   entity.ExchangeRateRepository
	provider   entity.ExchangeRateProvider
}

func NewCurrencyService(
	cryptoRepo entity.CryptoCurrencyRepository,
	rateRepo entity.ExchangeRateRepository,
	provider entity.ExchangeRateProvider) (*CurrencyService, error) {
	if cryptoRepo == nil {
		return nil, errors.New("cryptocurrency repository is nil")
	}
	if rateRepo == nil {
		return nil, errors.New("exchange rate repository is nil")
	}
	if provider == nil {
		return nil, errors.New("exchange rate provider is nil")
	}
	return &CurrencyService{cryptoRepo: cryptoRepo, rateRepo: rateRepo, provider: provider}, nil
}

// StartUpdate обновляет курсы сразу при вызове, затем каждые interval.
func (s *CurrencyService) StartUpdate(ctx context.Context, interval time.Duration) {
	_ = s.UpdateRates(ctx) // ошибку на старте игнорируем осознанно, следующий тик попробует снова

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = s.UpdateRates(ctx)
		}
	}
}

// resolveCryptocurrency добавление новой валюты
func (s *CurrencyService) resolveCryptocurrency(ctx context.Context, symbols []string) ([]entity.CryptoCurrency, error) {
	found, err := s.cryptoRepo.FindBySymbol(ctx, symbols)
	if err != nil {
		return nil, err
	}
	if len(found) == len(symbols) {
		return found, nil // Flow 1 для всех запрошенных символов сразу
	}

	alreadyFound := make(map[string]struct{}, len(found))
	for _, c := range found {
		alreadyFound[c.Symbol] = struct{}{}
	}

	// Flow 2: для каждого недостающего символа
	var newCurrencies []entity.CryptoCurrency
	priceBySymbol := make(map[string]float64)
	for _, symbol := range symbols {
		if _, ok := alreadyFound[symbol]; ok {
			continue
		}
		info, err := s.provider.GetCurrencyInfo(ctx, symbol)
		if err != nil {
			return nil, err
		}
		newCurrencies = append(newCurrencies, entity.CryptoCurrency{Symbol: info.Symbol, NameCurrency: info.Name})
		priceBySymbol[info.Symbol] = info.PriceUSD
	}
	if len(newCurrencies) == 0 {
		return found, nil
	}

	if err := s.cryptoRepo.Create(ctx, newCurrencies); err != nil {
		return nil, err
	}

	now := time.Now()
	initialRates := make([]entity.ExchangeRate, 0, len(newCurrencies))
	for _, c := range newCurrencies {
		initialRates = append(initialRates, entity.ExchangeRate{
			CryptoCurrencyID: c.ID,
			Price:            priceBySymbol[c.Symbol],
			Time:             now,
		})
	}
	if err := s.rateRepo.Save(ctx, initialRates); err != nil {
	}
	return append(found, newCurrencies...), nil
}

// UpdateRates обновить курсы
func (s *CurrencyService) UpdateRates(ctx context.Context) error {
	currencies, err := s.cryptoRepo.FindAll(ctx)
	if err != nil {
		return err
	}
	if len(currencies) == 0 {
		return nil
	}
	rates, err := s.provider.GetRates(ctx, currencies)
	if err != nil {
		return err
	}
	return s.rateRepo.Save(ctx, rates)
}

// GetCurrentRates возвращает последний курс по каждому запрошенному символу.
func (s *CurrencyService) GetCurrentRates(ctx context.Context, symbols []string) ([]dto.Rate, error) {
	currencies, err := s.resolveCryptocurrency(ctx, symbols)
	if err != nil {
		return nil, err
	}

	rates := make([]dto.Rate, 0, len(currencies))
	for _, currency := range currencies {
		rate, err := s.rateRepo.GetLatest(ctx, currency.ID)
		if err != nil {
			return nil, err
		}
		rates = append(rates, dto.Rate{Symbol: currency.Symbol, PriceUSD: rate.Price, UpdatedAt: rate.Time})
	}
	return rates, nil
}

// GetStatistics возвращает текущую цену, min/max за 24ч и % изменения ЗА ПОСЛЕДНИЙ ЧАС
func (s *CurrencyService) GetStatistics(ctx context.Context, symbols []string) ([]dto.Statistics, error) {
	currencies, err := s.resolveCryptocurrency(ctx, symbols)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	result := make([]dto.Statistics, 0, len(currencies))

	for _, currency := range currencies {
		current, err := s.rateRepo.GetLatest(ctx, currency.ID)
		if err != nil {
			return nil, err
		}

		min, max, err := s.rateRepo.GetMinMaxSince(ctx, currency.ID, now.Add(-24*time.Hour))
		if err != nil {
			return nil, err
		}

		var changePercent float64
		hourAgo, err := s.rateRepo.GetClosestBefore(ctx, currency.ID, now.Add(-1*time.Hour))
		if err == nil && hourAgo != nil && hourAgo.Price > 0 {
			changePercent = (current.Price - hourAgo.Price) / hourAgo.Price * 100
		}
		// если hourAgo == nil (данных за последний час ещё нет) — ChangePercent остаётся 0,
		// это ожидаемо для только что зарегистрированной монеты

		result = append(result, dto.Statistics{
			Symbol:        currency.Symbol,
			CurrentPrice:  current.Price,
			MinPrice:      min,
			MaxPrice:      max,
			ChangePercent: changePercent,
		})
	}
	return result, nil
}
