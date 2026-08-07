package entity

import "context"

type CryptoCurrencyRepository interface {
	FindAll(ctx context.Context) ([]CryptoCurrency, error)

	// FindBySymbol возвращает те валюты из запрошенных символов, что уже есть в БД.
	FindBySymbol(ctx context.Context, symbols []string) ([]CryptoCurrency, error)

	// Create сохраняет пачку новых валют и проставляет сгенерированные ID
	Create(ctx context.Context, currencies []CryptoCurrency) error
}
