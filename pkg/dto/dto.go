package dto

import (
	"time"
)

type CoinDTO struct {
	Title        string    `json:"title"`
	Price        float64   `json:"price"`
	CreationTime time.Time `json:"creationTime"`
}
type CoinsDTO []CoinDTO
