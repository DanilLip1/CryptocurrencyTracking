package entity

import (
	"time"

	"github.com/pkg/errors"
)

type Coin struct {
	Title        string
	Price        float64
	CreationTime time.Time
}

func NewCoin(title string, price float64, creationTime time.Time) (*Coin, error) {
	if len(title) == 0 {
		return nil, errors.Wrap(ErrInvalidParams, "title is required")
	}
	if price <= 0 {
		return nil, errors.Wrap(ErrInvalidParams, "price must be greater than zero")
	}
	if creationTime.IsZero() {
		return nil, errors.Wrap(ErrInvalidParams, "creation time is required")
	}
	return &Coin{Title: title, Price: price, CreationTime: creationTime}, nil
}
