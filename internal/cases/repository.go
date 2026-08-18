package cases

import (
	"context"
	"cryptocurrency/internal/entity"
)

type Repository interface {
	SaveCoinPrices(ctx context.Context, coins []entity.Coin) error

	GetLatestPrices(ctx context.Context, titles []string) ([]entity.Coin, error)

	GetMinPrices(ctx context.Context, titles []string) ([]entity.Coin, error)

	GetMaxPrices(ctx context.Context, titles []string) ([]entity.Coin, error)

	GetPriceChangePercent(ctx context.Context, titles []string) ([]entity.Coin, error)

	GetTitles(ctx context.Context) ([]string, error)

	AddTrackedTitles(ctx context.Context, titles []string) error
}
