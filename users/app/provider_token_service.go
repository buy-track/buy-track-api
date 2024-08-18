package app

import (
	"my-stocks/common/errors"
	"my-stocks/domains"
	"my-stocks/users/persistance/repositories"
)

type ProviderTokenService struct {
	providerReader repositories.ProviderTokenReader
	providerWriter repositories.ProviderTokenWriter
}

func NewProviderTokenService(providerReader repositories.ProviderTokenReader, providerWriter repositories.ProviderTokenWriter) *ProviderTokenService {
	return &ProviderTokenService{providerReader: providerReader, providerWriter: providerWriter}
}

func (a *ProviderTokenService) AddProviderToken(userId string, token string, provider domains.Provider) error {
	if a.ExistsProviderToken(token, provider) {
		return errors.DuplicateError{Data: token}
	}
	_, err := a.providerWriter.Create(domains.ProviderToken{
		UserId:     userId,
		ProviderId: token,
		Provider:   provider,
	})
	return err
}

func (a *ProviderTokenService) ExistsProviderToken(token string, provider domains.Provider) bool {
	_, err := a.providerReader.Get(token, provider)
	if err != nil {
		return false
	}
	return true
}

func (a *ProviderTokenService) DeleteAllProviderToken(userId string) bool {
	err := a.providerWriter.DeleteAllByUserId(userId)
	if err != nil {
		return false
	}
	return true
}
