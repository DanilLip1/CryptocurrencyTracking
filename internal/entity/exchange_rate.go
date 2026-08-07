package entity

import "time"

// ExchangeRate — курс конкретной валюты в конкретный момент времени
type ExchangeRate struct {
	CryptoCurrencyID int
	Time             time.Time
	Price            float64
}
