package postgres

import (
	"gorm.io/gorm"
	"my-stocks/coins/persistance/repositories"
	"my-stocks/common/errors"
	"my-stocks/domains"
)

// CoinRepository starts
type CoinRepository struct {
	db *gorm.DB
}

func NewCoinRepository(db *gorm.DB) repositories.CoinProvider {
	return &CoinRepository{db: db}
}

func (c CoinRepository) List(query repositories.CoinQueryable) domains.ListItem[*domains.Coin] {
	var tmp CoinList
	if query.Offset == 0 {
		query.Offset = 15
	}

	c.db.Where("status = ?", domains.Active).Limit(query.Limit).Offset(query.Offset - 1).Find(&tmp)
	return domains.ListItem[*domains.Coin]{
		Data:   tmp.ToCoinsDomain(),
		Limit:  query.Limit,
		Offset: query.Offset,
	}
}

func (c CoinRepository) All(symbols ...string) []*domains.Coin {
	var tmp CoinList
	query := c.db.Where("status = ?", domains.Active)
	if len(symbols) != 0 {
		query.Where("symbol in (?)", symbols)
	}
	query.Order("id desc").Find(&tmp)
	return tmp.ToCoinsDomain()
}

func (c CoinRepository) GetBySymbol(symbol string) (*domains.Coin, error) {
	var tmp CoinEntity
	c.db.Where("symbol = ?", symbol).First(&tmp)
	if tmp.ID != 0 {
		return tmp.ToCoin(), nil
	}
	return nil, errors.NotFoundError{}
}

func (c CoinRepository) GetById(id string) (*domains.Coin, error) {
	var tmp CoinEntity
	c.db.Where("id = ?", id).First(&tmp)
	if tmp.ID != 0 {
		return tmp.ToCoin(), nil
	}
	return nil, errors.NotFoundError{}
}

func (c CoinRepository) GetByIds(id []string) []*domains.Coin {
	var tmp CoinList
	c.db.Where("status = ?", domains.Active).Where("id in (?)", id).Order("id desc").Find(&tmp)
	return tmp.ToCoinsDomain()
}
