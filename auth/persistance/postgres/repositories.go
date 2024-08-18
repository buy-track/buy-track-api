package postgres

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"my-stocks/auth/persistance/repositories"
	"my-stocks/common/errors"
	"my-stocks/domains"
	"strconv"
	"time"
)

type AccessTokenRepository struct {
	db *gorm.DB
}

func NewAccessTokenRepository(db *gorm.DB) repositories.AccessTokenProvider {
	return &AccessTokenRepository{db: db}
}

func (a AccessTokenRepository) Get(token string) (*domains.Token, error) {
	var tmp AccessTokenEntity
	a.db.Where("token = ?", token).First(&tmp)
	if tmp.ID != 0 {
		return tmp.ToToken(), nil
	}
	return nil, errors.NotFoundError{}
}

func (a AccessTokenRepository) Create(token domains.Token) (*domains.Token, error) {
	userId, err := strconv.ParseUint(token.UserId, 10, 64)
	if err != nil {
		log.Errorln("Error during convert string to uint64", err)
		return nil, err
	}

	tmp := AccessTokenEntity{
		ExpiredAt: token.ExpiredAt,
		UserId:    userId,
		Token:     token.Token,
		Dates: Dates{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	a.db.Create(&tmp)
	return tmp.ToToken(), nil
}

func (a AccessTokenRepository) Update(token domains.Token) (*domains.Token, error) {
	tmp := AccessTokenEntity{
		ExpiredAt: token.ExpiredAt,
	}
	a.db.Where("token = ?", token.Token).Save(&tmp)
	return tmp.ToToken(), nil
}

func (a AccessTokenRepository) Delete(token domains.Token) error {
	a.db.Where("token = ?", token.Token).Delete(&AccessTokenEntity{})
	return nil
}

func (a AccessTokenRepository) DeleteAllByUserId(userId string) error {
	a.db.Where("user_id = ?", userId).Delete(&AccessTokenEntity{})
	return nil
}
