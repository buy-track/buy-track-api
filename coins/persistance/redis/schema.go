package redisRepository

import (
	"my-stocks/domains"
)

type CoinPriceList []*CoinPriceEntity

func (c CoinPriceList) ToCoinsPriceDomain() []*domains.CoinPrice {
	tmp := make([]*domains.CoinPrice, len(c))
	for i, item := range c {
		tmp[i] = item.ToCoinPrice()
	}
	return tmp
}

type CoinPriceEntity struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
	CoinId    string  `json:"coin_id"`
}

func (coin *CoinPriceEntity) ToCoinPrice() *domains.CoinPrice {
	return &domains.CoinPrice{
		Timestamp: coin.Timestamp,
		Price:     coin.Price,
		CoinId:    coin.CoinId,
	}
}
