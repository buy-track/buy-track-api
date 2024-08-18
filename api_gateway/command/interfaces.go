package command

import "my-stocks/domains"

type AuthCommander interface {
	Login(userId string) (*domains.Token, error)
	Logout(token string) error
	RevokeAllTokens(userId string) error
}

type UserCommander interface {
	Create(data *domains.User) (*domains.User, error)
	AddProviderToken(userId string, token string, provider domains.Provider) error
}

type CoinCommander interface {
}
