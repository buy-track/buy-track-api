package domains

import "time"

type Token struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt *time.Time
	UserId    string
	Token     string
}

func NewToken(expiredAt *time.Time, userId string, token string, createdAt time.Time, updatedAt time.Time) *Token {
	return &Token{CreatedAt: createdAt, UpdatedAt: updatedAt, ExpiredAt: expiredAt, UserId: userId, Token: token}
}

type Provider uint8

const (
	Google Provider = 1
	Apple  Provider = 2
)

type ProviderToken struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserId     string
	ProviderId string
	Provider   Provider
}
