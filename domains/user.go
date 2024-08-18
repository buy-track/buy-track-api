package domains

import (
	"encoding/json"
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Password  string
	Email     string
	Name      string
	ID        string
}

func NewUser(createdAt time.Time, updatedAt time.Time, email string, name string, ID string) *User {
	return &User{CreatedAt: createdAt, UpdatedAt: updatedAt, Email: email, Name: name, ID: ID}
}

func (user User) ToJson() ([]byte, error) {
	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
