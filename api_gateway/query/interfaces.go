package query

import (
	"my-stocks/domains"
)

type AuthQuery interface {
	VerifyToken(token string) (string, error)
}

type UserQuery interface {
	EmailExists(email string) bool
	GetByEmail(email string) (*domains.User, error)
	GetById(id string) (*domains.User, error)
	CheckPassword(email string, password string) (*domains.User, error)
}

type CoinQuery interface {
	Paginate(limit, offset int64) domains.ListItem[*domains.Coin]
	GetById(id string) (*domains.Coin, error)
	GetByIds(id []string) []*domains.Coin
	GetBySymbol(symbol string) (*domains.Coin, error)
}
