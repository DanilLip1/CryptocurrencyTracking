package entity

import (
	"errors"
	"time"
)

type Coin struct {
	Title        string
	Price        float64
	CreationTime time.Time
}

func NewCoin(title string, price float64, creationTime time.Time) (*Coin, error) {
	if len(title) == 0 {
		return nil, errors.New("title is required")
	}
	if price <= 0 {
		return nil, errors.New("value is required")
	}
	if creationTime.IsZero() {
		return nil, errors.New("creationTime is required")
	}
	return &Coin{Title: title, Price: price, CreationTime: creationTime}, nil
}
