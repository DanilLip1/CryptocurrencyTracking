package public

import (
	"context"
	"cryptocurrency/internal/entity"
)

type Service interface {
	GetLatestPrices(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetMinPrices(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetMaxPrices(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetPriceChangePercent(ctx context.Context, titles []string) ([]entity.Coin, error)
}
