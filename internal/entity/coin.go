package entity

import (
	"errors"
	"time"
)

type Coin struct {
	Title        string
	Value        float64
	CreationTime time.Time
}

func NewCoin(title string, value float64, creationTime time.Time) (*Coin, error) {
	if len(title) == 0 {
		return nil, errors.New("title is required")
	}
	if value <= 0 {
		return nil, errors.New("value is required")
	}
	if creationTime.IsZero() {
		return nil, errors.New("creationTime is required")
	}
	return &Coin{Title: title, Value: value, CreationTime: creationTime}, nil
}
