package http

import (
	"cryptocurrency/internal/entity"
	"time"
)

type CoinDTO struct {
	Title        string    `json:"title"`
	Price        float64   `json:"price"`
	CreationTime time.Time `json:"creationTime"`
}

func NewCoinDTO(coin entity.Coin) CoinDTO {
	return CoinDTO{
		Title:        coin.Title,
		Price:        coin.Price,
		CreationTime: coin.CreationTime,
	}
}
