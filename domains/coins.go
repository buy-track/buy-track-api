package domains

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type CoinType int

const (
	Unknown CoinType = iota
	Fiat
	Crypto
)

type CoinStatus int

const (
	InActive CoinStatus = iota
	Active
)

func (b *CoinType) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *CoinType) UnmarshalText(text []byte) error {
	*b = CoinTypeFromText(string(text))
	return nil
}

func (b *CoinType) String() string {
	text, err := b.CoinTypeToText()
	if err != nil {
		return text
	}
	return text
}

func CoinTypeFromText(text string) CoinType {
	switch strings.ToLower(text) {
	default:
		return Unknown
	case "fiat":
		return Fiat
	case "crypto":
		return Crypto
	}
}

func (b *CoinType) CoinTypeToText() (string, error) {
	switch *b {
	case Fiat:
		return "fiat", nil
	case Crypto:
		return "crypto", nil
	case Unknown:
		fallthrough
	default:
		return "unknown", errors.New("unknown type")
	}
}

type Coin struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Type      CoinType  `json:"type"`
	Price     float64   `json:"price"`
	Symbol    string    `json:"symbol"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon,omitempty"`
	ID        string    `json:"id"`
}

func NewCoin(ID string, name string, symbol string, icon string, createdAt time.Time, updatedAt time.Time) *Coin {
	return &Coin{CreatedAt: createdAt, UpdatedAt: updatedAt, Symbol: symbol, Name: name, Icon: icon, ID: ID}
}

func NewCoinFromJson(data []byte) (*Coin, error) {
	var tmp Coin
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, err
	}
	return &tmp, nil
}

func (coin Coin) ToJson() ([]byte, error) {
	bytes, err := json.Marshal(coin)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type CoinPrice struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
	CoinId    string  `json:"coin_id"`
}
