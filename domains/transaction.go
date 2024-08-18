package domains

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type TransactionType int8

const (
	UnknownTransactionType TransactionType = 0
	Decrease               TransactionType = -1
	Increase               TransactionType = 1
)

func (t *TransactionType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *TransactionType) UnmarshalText(text []byte) error {
	*t = TransactionTypeFromText(string(text))
	return nil
}

func (t *TransactionType) String() string {
	text, err := t.TransactionTypeToText()
	if err != nil {
		return text
	}
	return text
}

func TransactionTypeFromText(text string) TransactionType {
	switch strings.ToLower(text) {
	case "decrease":
		return Decrease
	case "increase":
		return Increase
	default:
		return UnknownTransactionType
	}
}

func (t *TransactionType) TransactionTypeToText() (string, error) {
	switch *t {
	case Increase:
		return "increase", nil
	case Decrease:
		return "decrease", nil
	case UnknownTransactionType:
		fallthrough
	default:
		return "unknown", errors.New("unknown type")
	}
}

type Transaction struct {
	TrackingCode string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ParentId     string
	Reason       string
	UserId       string
	CoinId       string
	Amount       float64
	Type         TransactionType
	ID           string
}

func NewTransaction(trackingCode string, createdAt time.Time, updatedAt time.Time, parentId string, reason string, userId string, coinId string, amount float64, Type TransactionType, ID string) *Transaction {
	return &Transaction{TrackingCode: trackingCode, CreatedAt: createdAt, UpdatedAt: updatedAt, ParentId: parentId, Reason: reason, UserId: userId, CoinId: coinId, Amount: amount, Type: Type, ID: ID}
}

func (transaction Transaction) ToJson() ([]byte, error) {
	bytes, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
