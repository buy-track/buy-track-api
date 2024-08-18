package app

import (
	"fmt"
	"my-stocks/api-gateway/command"
	"my-stocks/api-gateway/query"
	"my-stocks/domains"
)

type CoinService struct {
	command command.CoinCommander
	query   query.CoinQuery
}

func NewCoinService(command command.CoinCommander, query query.CoinQuery) *CoinService {
	return &CoinService{command: command, query: query}
}

func (c CoinService) PaginateList(limit, offset int64) domains.ListItem[*domains.Coin] {
	fmt.Println("11xxxxxx")
	return c.query.Paginate(limit, offset)
}

func (c CoinService) GetById(id string) (*domains.Coin, error) {
	return c.query.GetById(id)
}

func (c CoinService) GetByIds(id []string) []*domains.Coin {
	return c.query.GetByIds(id)
}

func (c CoinService) GetBySymbol(symbol string) (*domains.Coin, error) {
	return c.query.GetBySymbol(symbol)
}
