package entity

import "context"

// CurrencyInfo — то, что провайдер знает о монете
type CurrencyInfo struct {
	Symbol   string
	Name     string
	PriceUSD float64
}

type ExchangeRateProvider interface {
	// GetRates —запрос актуальных цен для уже отслеживаемых валют.
	GetRates(ctx context.Context, currencies []CryptoCurrency) ([]ExchangeRate, error)

	// GetCurrencyInfo — точечный запрос при первом обращении к неизвестному символу.
	GetCurrencyInfo(ctx context.Context, symbol string) (CurrencyInfo, error)
}
