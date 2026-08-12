package cases

import (
	"context"
	"cryptocurrency/internal/entity"
	"time"
)

type Repository interface {

	//SaveRates сохраняет
	SaveRates(ctx context.Context, rates []entity.Coin) error

	// GetLast — самая свежая котировка конкретной валюты
	GetLast(ctx context.Context, title []string) (*[]entity.Coin, error)

	//GetMin24h - минимальная цена 24h
	GetMin24h(ctx context.Context, title []string) (float64, error)

	// GetMax24h — максимальная цена 24h
	GetMax24h(ctx context.Context, title []string) (float64, error)

	// GetPercent — процент
	GetPercent(ctx context.Context, title []string, at time.Time) (*[]entity.Coin, error)

	// GetAll возвращает все отслеживаемые валюты
	GetAll(ctx context.Context) (*[]entity.Coin, error)

	// CreateCoin сохраняет пачку новых валют
	CreateCoin(ctx context.Context, coin []entity.Coin) error
}
