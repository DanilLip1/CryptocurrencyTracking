package cases

import (
	"context"
	"cryptocurrency/internal/entity"
)

type Repository interface {

	//SaveRates сохраняет существующие
	SaveRates(ctx context.Context, rates []entity.Coin) error

	// GetLatestRates — самая свежая котировка конкретной валюты
	GetLatestRates(ctx context.Context, titles []string) ([]entity.Coin, error)

	//GetMinPrices - минимальная цена
	GetMinPrices(ctx context.Context, titles []string) ([]entity.Coin, error)

	// GetMaxPrices — максимальная цена
	GetMaxPrices(ctx context.Context, titles []string) ([]entity.Coin, error)

	// GetPriceChangePercent — процент
	GetPriceChangePercent(ctx context.Context, titles []string) ([]entity.Coin, error)

	// GetAllCoins возвращает все отслеживаемые валюты
	GetAllCoins(ctx context.Context) ([]string, error)

	// CreateCoins добавляет новые монеты
	CreateCoins(ctx context.Context, coins []entity.Coin) error
}
