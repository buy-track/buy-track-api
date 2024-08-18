package repositories

import "my-stocks/domains"

type UserReader interface {
	GetByEmail(email string) (*domains.User, error)
	EmailExists(email string) bool
	GetById(id string) (*domains.User, error)
}

type UserWriter interface {
	Create(user domains.User) (*domains.User, error)
}

type UserProvider interface {
	UserReader
	UserWriter
}
