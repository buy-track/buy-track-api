package app

import (
	"my-stocks/api-gateway/command"
	"my-stocks/api-gateway/query"
	"my-stocks/domains"
)

type AuthService struct {
	command command.AuthCommander
	query   query.AuthQuery
}

func NewAuthService(command command.AuthCommander, query query.AuthQuery) *AuthService {
	return &AuthService{command: command, query: query}
}

func (a AuthService) Login(userId string) (*domains.Token, error) {
	return a.command.Login(userId)
}

func (a AuthService) Logout(token string) error {
	return a.command.Logout(token)
}

func (a AuthService) RevokeAllTokens(userId string) error {
	return a.command.RevokeAllTokens(userId)
}

func (a AuthService) VerifyToken(token string) (string, error) {
	return a.query.VerifyToken(token)
}
