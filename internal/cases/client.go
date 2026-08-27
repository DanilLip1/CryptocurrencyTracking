package cases

import (
	"context"
	"cryptocurrency/internal/entity"
)

type Provider interface {
	GetRates(ctx context.Context, titles []string) ([]entity.Coin, error)
}
