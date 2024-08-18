package postgres

import (
	"my-stocks/domains"
	"strconv"
	"time"
)

type Identifier struct {
	ID uint64 `gorm:"primaryKey"`
}

type Dates struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:nano"`
}

type SoftDelete struct {
	DeletedAt *time.Time `gorm:"type:timestamp;"`
}

type UserEntity struct {
	Password string `gorm:"size:144;"`
	Email    string `gorm:"size:144;index:email_index,unique;not null"`
	Name     string `gorm:"size:144;not null"`
	Identifier
	Dates
}

func (user *UserEntity) ToUser() *domains.User {
	return &domains.User{
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Password:  user.Password,
		Email:     user.Email,
		Name:      user.Name,
		ID:        strconv.FormatUint(user.ID, 10),
	}
}

type ProviderTokenEntity struct {
	UserId     uint64 `gorm:"index:user_id_index;not null"`
	ProviderId string `gorm:"size:144;index:provider_index,unique;not null"`
	Provider   uint8  `gorm:"index:provider_index,unique;not null"`
	Identifier
	Dates
}

func (e *ProviderTokenEntity) ToProviderToken() *domains.ProviderToken {
	return &domains.ProviderToken{
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
		UserId:     strconv.FormatUint(e.UserId, 10),
		ProviderId: e.ProviderId,
		Provider:   domains.Provider(e.Provider),
	}
}
