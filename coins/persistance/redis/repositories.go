package redisRepository

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"my-stocks/coins/persistance/repositories"
	"my-stocks/domains"
)

// CoinPriceRepository starts
type CoinPriceRepository struct {
	db     *redis.Client
	ctx    context.Context
	prefix string
}

func NewCoinPriceRepository(db *redis.Client, ctx context.Context, prefix string) repositories.CoinPriceProvider {
	if prefix == "" {
		prefix = "coins:prices:"
	} else {
		prefix = prefix + ":prices:"
	}

	return &CoinPriceRepository{db: db, ctx: ctx, prefix: prefix}
}

func (c CoinPriceRepository) ListByCoinIds(ids []string) []*domains.CoinPrice {
	pipe := c.db.Pipeline()
	for _, id := range ids {
		pipe.LIndex(c.ctx, c.generateKey(id), 0)
	}

	cmders, err := pipe.Exec(c.ctx)
	if err != nil && err != redis.Nil {
		return []*domains.CoinPrice{}
	}

	tmp := make(CoinPriceList, 0)
	for _, cmd := range cmders {
		val, err := cmd.(*redis.StringCmd).Result()
		if err == nil && err != redis.Nil {
			var t CoinPriceEntity
			json.Unmarshal([]byte(val), &t)
			tmp = append(tmp, &t)
		}
	}
	return tmp.ToCoinsPriceDomain()
}

func (c CoinPriceRepository) GetByCoinId(id string) (*domains.CoinPrice, error) {
	result, err := c.db.LIndex(c.ctx, c.generateKey(id), 0).Result()
	if err != nil {
		return nil, err
	}
	var t CoinPriceEntity
	err = json.Unmarshal([]byte(result), &t)

	if err != nil {
		return nil, err
	}
	return t.ToCoinPrice(), err
}

func (c CoinPriceRepository) UpdateOrCreate(price domains.CoinPrice) bool {
	marshal, err := json.Marshal(price)
	if err != nil {
		return false
	}
	c.db.LPush(c.ctx, c.generateKey(price.CoinId), string(marshal))
	c.db.LTrim(c.ctx, c.generateKey(price.CoinId), 0, 0)

	return true
}

func (c CoinPriceRepository) generateKey(key string) string {
	return c.prefix + key
}
