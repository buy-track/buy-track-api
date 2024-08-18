package postgres

import (
	"my-stocks/domains"
	"strconv"
	"time"
)

type Identifier struct {
	ID uint64 `gorm:"primaryKey"`
}

type Dates struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:nano"`
}

type SoftDelete struct {
	DeletedAt *time.Time `gorm:"type:timestamp;"`
}

type CoinList []*CoinEntity

func (c CoinList) ToCoinsDomain() []*domains.Coin {
	tmp := make([]*domains.Coin, len(c))
	for i, item := range c {
		tmp[i] = item.ToCoin()
	}
	return tmp
}

type CoinEntity struct {
	Status int    `gorm:"default:1"`
	Symbol string `gorm:"size:144;index:symbol_index,unique;not null"`
	Type   int    `gorm:"default:0"`
	Icon   string `gorm:"size:144"`
	Name   string `gorm:"size:144;not null"`
	Identifier
	Dates
}

func (coin *CoinEntity) ToCoin() *domains.Coin {
	return &domains.Coin{
		CreatedAt: coin.CreatedAt,
		UpdatedAt: coin.UpdatedAt,
		Type:      domains.CoinType(coin.Type),
		Symbol:    coin.Symbol,
		Name:      coin.Name,
		Icon:      coin.Icon,
		ID:        strconv.FormatUint(coin.ID, 10),
	}
}
