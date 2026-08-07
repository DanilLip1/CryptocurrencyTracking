package entity

import (
	"context"
	"time"
)

// ExchangeRateRepository — контракт хранения истории котировок.
type ExchangeRateRepository interface {
	// Save сохраняет
	Save(ctx context.Context, rates []ExchangeRate) error

	// GetLatest — самая свежая котировка конкретной валюты.
	GetLatest(ctx context.Context, cryptoCurrencyID int) (ExchangeRate, error)

	// GetMinMaxSince — минимум/максимум цены с момента since по сейчас.
	GetMinMaxSince(ctx context.Context, cryptoCurrencyID int, since time.Time) (min, max float64, err error)

	// GetClosestBefore — котировка, ближайшая по времени к at, но не позже него.
	// Нужна для расчёта %
	GetClosestBefore(ctx context.Context, cryptoCurrencyID int, at time.Time) (*ExchangeRate, error)
}
