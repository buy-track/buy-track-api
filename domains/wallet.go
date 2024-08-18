package domains

import (
	"encoding/json"
	"time"
)

type Wallet struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    string
	CoinId    string
	Amount    float64
	ID        string
}

func NewWallet(createdAt time.Time, updatedAt time.Time, userId string, coinId string, amount float64, ID string) *Wallet {
	return &Wallet{CreatedAt: createdAt, UpdatedAt: updatedAt, UserId: userId, CoinId: coinId, Amount: amount, ID: ID}
}

func (wallet Wallet) ToJson() ([]byte, error) {
	bytes, err := json.Marshal(wallet)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
