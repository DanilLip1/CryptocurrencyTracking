package cases

import (
	"context"
	"cryptocurrency/internal/entity"
)

type ExchangeRateProvider interface {
	// GetRates возвращает актуальные курсы
	GetRates(ctx context.Context, title []string) ([]entity.Coin, error)
}
