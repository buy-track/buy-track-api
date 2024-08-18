package repositories

import "my-stocks/domains"

type CoinReader interface {
	List(query CoinQueryable) domains.ListItem[*domains.Coin]
	All(symbols ...string) []*domains.Coin
	GetBySymbol(symbol string) (*domains.Coin, error)
	GetById(id string) (*domains.Coin, error)
	GetByIds(id []string) []*domains.Coin
}

type CoinWriter interface {
}

type CoinProvider interface {
	CoinReader
	CoinWriter
}

type CoinQueryable struct {
	Limit  int
	Offset int
}

type CoinPriceReader interface {
	ListByCoinIds(ids []string) []*domains.CoinPrice
	GetByCoinId(id string) (*domains.CoinPrice, error)
}

type CoinPriceWriter interface {
	UpdateOrCreate(price domains.CoinPrice) bool
}

type CoinPriceProvider interface {
	CoinPriceReader
	CoinPriceWriter
}
