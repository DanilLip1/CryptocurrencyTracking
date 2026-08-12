package cases

import (
	"context"
	"cryptocurrency/internal/entity"
)

type Repository interface {

	//SaveRates сохраняет
	SaveRates(ctx context.Context, rates []entity.Coin) error

	// GetLast — самая свежая котировка конкретной валюты
	GetLast(ctx context.Context, title []string) ([]entity.Coin, error)

	//GetMinPrice - минимальная цена 24h
	GetMinPrice(ctx context.Context, title []string) ([]entity.Coin, error)

	// GetMaxPrice — максимальная цена 24h
	GetMaxPrice(ctx context.Context, title []string) ([]entity.Coin, error)

	// GetChangePercent — процент
	GetChangePercent(ctx context.Context, title []string) ([]entity.Coin, error)

	// GetAll возвращает все отслеживаемые валюты
	GetAll(ctx context.Context) ([]string, error)

	// CreateCoin сохраняет пачку новых валют
	CreateCoin(ctx context.Context, coin []entity.Coin) error
}
