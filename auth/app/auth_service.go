package app

import (
	"crypto/rand"
	"encoding/base64"
	"my-stocks/auth/persistance/repositories"
	"my-stocks/domains"
	"strconv"
	"time"
)

type AuthService struct {
	accessTokenReader repositories.AccessTokenReader
	accessTokenWriter repositories.AccessTokenWriter
}

func NewAuthService(accessTokenReader repositories.AccessTokenReader, accessTokenWriter repositories.AccessTokenWriter) *AuthService {
	return &AuthService{accessTokenReader: accessTokenReader, accessTokenWriter: accessTokenWriter}
}

func (a *AuthService) GenerateAccessToken(userId string) (*domains.Token, error) {
	token := a.generateUniqueToken()
	expiredAt := time.Now().AddDate(0, 1, 0)
	tmp := domains.Token{
		ExpiredAt: &expiredAt,
		UserId:    userId,
		Token:     token,
	}
	created, err := a.accessTokenWriter.Create(tmp)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (a *AuthService) VerifyAccessToken(token string) (*domains.Token, error) {
	return a.accessTokenReader.Get(token)
}

func (a *AuthService) DeleteAccessToken(token string) error {
	return a.accessTokenWriter.Delete(domains.Token{
		Token: token,
	})
}

func (a *AuthService) DeleteAllAccessToken(userId string) error {
	return a.accessTokenWriter.DeleteAllByUserId(userId)
}

func (a *AuthService) generateUniqueToken() string {
	tokenBytes := make([]byte, 30)
	_, err := rand.Read(tokenBytes)

	if err != nil {
		return ""
	}
	str := append(tokenBytes, []byte(strconv.FormatInt(time.Now().UnixMicro(), 10))...)
	return base64.URLEncoding.EncodeToString(str)
}
