package repositories

import (
	"my-stocks/domains"
)

type AccessTokenReader interface {
	Get(token string) (*domains.Token, error)
}

type AccessTokenWriter interface {
	Create(token domains.Token) (*domains.Token, error)
	Update(token domains.Token) (*domains.Token, error)
	Delete(token domains.Token) error
	DeleteAllByUserId(userId string) error
}

type AccessTokenProvider interface {
	AccessTokenReader
	AccessTokenWriter
}
