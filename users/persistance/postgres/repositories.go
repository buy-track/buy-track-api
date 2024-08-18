package postgres

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"my-stocks/common/errors"
	"my-stocks/domains"
	"my-stocks/users/persistance/repositories"
	"strconv"
)

// UserRepository starts
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserProvider {
	return &UserRepository{db: db}
}

func (u UserRepository) GetByEmail(email string) (*domains.User, error) {
	var tmp UserEntity
	u.db.Where("email = ?", email).First(&tmp)
	if tmp.ID != 0 {
		return tmp.ToUser(), nil
	}
	return nil, errors.NotFoundError{}
}

func (u UserRepository) EmailExists(email string) bool {
	var tmp UserEntity
	u.db.Where("email = ?", email).First(&tmp)
	if tmp.ID != 0 {
		return true
	}
	return false
}

func (u UserRepository) GetById(id string) (*domains.User, error) {
	var tmp UserEntity
	u.db.Where("id = ?", id).First(&tmp)
	if tmp.ID != 0 {
		return tmp.ToUser(), nil
	}
	return nil, errors.NotFoundError{}
}

func (u UserRepository) Create(user domains.User) (*domains.User, error) {
	tmp := UserEntity{
		Password: user.Password,
		Email:    user.Email,
		Name:     user.Name,
	}
	u.db.Create(&tmp)
	return tmp.ToUser(), nil
}

// ProviderTokenRepository starts
type ProviderTokenRepository struct {
	db *gorm.DB
}

func NewProviderTokenRepository(db *gorm.DB) repositories.ProviderTokenProvider {
	return &ProviderTokenRepository{db: db}
}

func (p ProviderTokenRepository) Create(token domains.ProviderToken) (*domains.ProviderToken, error) {
	userId, err := strconv.ParseUint(token.UserId, 10, 64)
	if err != nil {
		log.Errorln("Error during convert string to uint64", err)
		return nil, err
	}

	tmp := ProviderTokenEntity{
		UserId:     userId,
		ProviderId: token.ProviderId,
		Provider:   uint8(token.Provider),
	}
	p.db.Create(&tmp)
	return tmp.ToProviderToken(), nil
}

func (p ProviderTokenRepository) Delete(token *domains.ProviderToken) error {
	p.db.Where("provider_id = ?", token.ProviderId).Delete(&ProviderTokenEntity{})
	return nil
}

func (p ProviderTokenRepository) DeleteAllByUserId(userId string) error {
	p.db.Where("user_id = ?", userId).Delete(&ProviderTokenEntity{})
	return nil
}

func (p ProviderTokenRepository) Get(token string, provider domains.Provider) (*domains.ProviderToken, error) {
	var tmp ProviderTokenEntity
	p.db.Where("provider_id = ? AND provider = ?", token, provider).First(&tmp)
	if tmp.ID != 0 {
		return tmp.ToProviderToken(), nil
	}
	return nil, errors.NotFoundError{}
}
