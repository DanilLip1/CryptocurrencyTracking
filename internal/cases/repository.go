package cases

import (
	"context"
	"cryptocurrency/internal/entity"
)

type Repository interface {
	SaveCoinPrices(ctx context.Context, coins []entity.Coin) error

	Get(ctx context.Context, title []string, opts ...Option) ([]entity.Coin, error)

	GetTitles(ctx context.Context) ([]string, error)

	AddTrackedTitles(ctx context.Context, titles []string) error
}
