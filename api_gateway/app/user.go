package app

import (
	"my-stocks/api-gateway/command"
	"my-stocks/api-gateway/query"
	"my-stocks/domains"
)

type UserService struct {
	command command.UserCommander
	query   query.UserQuery
}

func NewUserService(command command.UserCommander, query query.UserQuery) *UserService {
	return &UserService{command: command, query: query}
}

func (u UserService) Create(user *domains.User) (*domains.User, error) {
	return u.command.Create(user)
}

func (u UserService) AddProviderToken(userId string, token string, provider domains.Provider) error {
	return u.command.AddProviderToken(userId, token, provider)
}

func (u UserService) EmailExists(email string) bool {
	return u.query.EmailExists(email)
}

func (u UserService) GetByEmail(email string) (*domains.User, error) {
	return u.query.GetByEmail(email)
}

func (u UserService) GetById(id string) (*domains.User, error) {
	return u.query.GetById(id)
}

func (u UserService) CheckPassword(email, password string) (*domains.User, error) {
	return u.query.CheckPassword(email, password)
}
