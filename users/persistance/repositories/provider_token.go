package repositories

import (
	"my-stocks/domains"
)

type ProviderTokenReader interface {
	Get(token string, provider domains.Provider) (*domains.ProviderToken, error)
}

type ProviderTokenWriter interface {
	Create(token domains.ProviderToken) (*domains.ProviderToken, error)
	Delete(token *domains.ProviderToken) error
	DeleteAllByUserId(userId string) error
}

type ProviderTokenProvider interface {
	ProviderTokenWriter
	ProviderTokenReader
}
