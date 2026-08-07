package dto

import "time"

type Statistics struct {
	Symbol        string
	CurrentPrice  float64
	MinPrice      float64
	MaxPrice      float64
	ChangePercent float64
}

type Rate struct {
	Symbol    string
	PriceUSD  float64
	UpdatedAt time.Time
}
