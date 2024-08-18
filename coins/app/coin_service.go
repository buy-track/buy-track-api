package app

import (
	"my-stocks/coins/persistance/repositories"
	"my-stocks/domains"
)

type CoinService struct {
	coinReader  repositories.CoinReader
	coinWriter  repositories.CoinWriter
	priceWriter repositories.CoinPriceWriter
	priceReader repositories.CoinPriceReader
}

func NewCoinService(coinReader repositories.CoinReader, coinWriter repositories.CoinWriter, priceWriter repositories.CoinPriceWriter, priceReader repositories.CoinPriceReader) *CoinService {
	return &CoinService{coinReader: coinReader, coinWriter: coinWriter, priceWriter: priceWriter, priceReader: priceReader}
}

func (coin CoinService) GetAllActiveCoins(symbols ...string) []*domains.Coin {
	coins := coin.coinReader.All(symbols...)
	ids := make([]string, len(coins))
	for i, c := range coins {
		ids[i] = c.ID
	}
	prices := coin.priceReader.ListByCoinIds(ids)
	pricesHash := make(map[string]float64, len(prices))
	for _, price := range prices {
		pricesHash[price.CoinId] = price.Price
	}
	for _, c := range coins {
		c.Price = pricesHash[c.ID]
	}
	return coins
}

func (coin CoinService) Paginate(limit, offset int) domains.ListItem[*domains.Coin] {
	coins := coin.coinReader.List(repositories.CoinQueryable{
		Limit:  limit,
		Offset: offset,
	})
	ids := make([]string, len(coins.Data))
	for i, c := range coins.Data {
		ids[i] = c.ID
	}
	prices := coin.priceReader.ListByCoinIds(ids)
	pricesHash := make(map[string]float64, len(prices))
	for _, price := range prices {
		pricesHash[price.CoinId] = price.Price
	}

	for _, c := range coins.Data {
		c.Price = pricesHash[c.ID]
	}
	return coins
}

func (coin CoinService) GetById(id string) (*domains.Coin, error) {
	c, err := coin.coinReader.GetById(id)
	if err != nil {
		return nil, err
	}
	price, err := coin.priceReader.GetByCoinId(c.ID)
	if err != nil {
		return nil, err
	}
	c.Price = price.Price

	return c, nil
}

func (coin CoinService) GetBySymbol(symbol string) (*domains.Coin, error) {
	c, err := coin.coinReader.GetBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	price, err := coin.priceReader.GetByCoinId(c.ID)
	if err != nil {
		return nil, err
	}
	c.Price = price.Price
	return c, nil
}

func (coin CoinService) UpdatePrice(price domains.CoinPrice) bool {
	return coin.priceWriter.UpdateOrCreate(price)
}

func (coin CoinService) GetByIds(id []string) []*domains.Coin {
	coins := coin.coinReader.GetByIds(id)
	ids := make([]string, len(coins))
	for i, c := range coins {
		ids[i] = c.ID
	}
	prices := coin.priceReader.ListByCoinIds(ids)
	pricesHash := make(map[string]float64, len(prices))
	for _, price := range prices {
		pricesHash[price.CoinId] = price.Price
	}
	for _, c := range coins {
		c.Price = pricesHash[c.ID]
	}
	return coins
}
