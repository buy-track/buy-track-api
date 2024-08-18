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

type AccessTokenEntity struct {
	ExpiredAt *time.Time `gorm:"type:timestamp;"`
	UserId    uint64     `gorm:"index:user_id_index;not null"`
	Token     string     `gorm:"size:144;index:token_index,unique;not null"`
	Identifier
	Dates
}

func (access *AccessTokenEntity) ToToken() *domains.Token {
	return &domains.Token{
		CreatedAt: access.Dates.CreatedAt,
		UpdatedAt: access.Dates.UpdatedAt,
		ExpiredAt: access.ExpiredAt,
		UserId:    strconv.FormatUint(access.UserId, 10),
		Token:     access.Token,
	}
}
