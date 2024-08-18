package app

import (
	"golang.org/x/crypto/bcrypt"
	"my-stocks/domains"
	"my-stocks/users/persistance/repositories"
)

type CreateUserDto struct {
	Password string
	Email    string
	Name     string
}

type UserService struct {
	userReader repositories.UserReader
	userWriter repositories.UserWriter
}

func NewUserService(userReader repositories.UserReader, userWriter repositories.UserWriter) *UserService {
	return &UserService{userReader: userReader, userWriter: userWriter}
}

func (u UserService) Register(dto CreateUserDto) (*domains.User, error) {
	tmp := domains.User{
		Email: dto.Email,
		Name:  dto.Name,
	}
	if dto.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		tmp.Password = string(hashed)
	}

	created, err := u.userWriter.Create(tmp)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (u UserService) EmailExists(email string) bool {
	return u.userReader.EmailExists(email)
}

func (u UserService) FindByEmail(email string) (*domains.User, error) {
	return u.userReader.GetByEmail(email)
}

func (u UserService) FindById(id string) (*domains.User, error) {
	return u.userReader.GetById(id)
}
