package grpc

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"my-stocks/api-gateway/query"
	"my-stocks/common/grpc/services"
	"my-stocks/domains"
	"time"
)

type CoinInquirer struct {
	client services.CoinServiceClient
	ctx    context.Context
}

func (c CoinInquirer) Paginate(limit, offset int64) domains.ListItem[*domains.Coin] {
	log.Println("offset", offset)
	log.Println("limit", limit)
	coins, err := c.client.Paginate(c.ctx, &services.CoinQueryRequest{
		Limit: limit,
		Page:  offset,
	})
	if err != nil {
		return domains.ListItem[*domains.Coin]{}
	}
	tmp := make([]*domains.Coin, len(coins.Coins.Coins))
	for i, coin := range coins.Coins.Coins {
		createdAt, _ := time.Parse("2024-06-10", coin.CreatedAt)
		updatedAt, _ := time.Parse("2024-06-10", coin.UpdatedAt)
		tmp[i] = &domains.Coin{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Type:      domains.CoinType(coin.Type),
			Price:     float64(coin.Price),
			Symbol:    coin.Symbol,
			Name:      coin.Name,
			Icon:      coin.Icon,
			ID:        coin.Id,
		}
	}
	fmt.Println(tmp[0])
	return domains.ListItem[*domains.Coin]{
		Data:   tmp,
		Limit:  int(coins.Limit),
		Offset: int(coins.Page),
	}
}

func (c CoinInquirer) GetById(id string) (*domains.Coin, error) {
	coin, err := c.client.GetById(c.ctx, &services.IdRequest{Id: id})
	if err != nil {
		return nil, err
	}
	createdAt, _ := time.Parse("2024-06-10", coin.CreatedAt)
	updatedAt, _ := time.Parse("2024-06-10", coin.UpdatedAt)
	return domains.NewCoin(coin.Id, coin.Name, coin.Symbol, coin.Icon, createdAt, updatedAt), nil
}

func (c CoinInquirer) GetByIds(id []string) []*domains.Coin {
	coins, err := c.client.GetByIds(c.ctx, &services.IdsRequest{Id: id})
	if err != nil {
		return []*domains.Coin{}
	}
	tmp := make([]*domains.Coin, len(coins.Coins))
	for i, coin := range coins.Coins {
		createdAt, _ := time.Parse("2024-06-10", coin.CreatedAt)
		updatedAt, _ := time.Parse("2024-06-10", coin.UpdatedAt)
		tmp[i] = domains.NewCoin(coin.Id, coin.Name, coin.Symbol, coin.Icon, createdAt, updatedAt)
	}
	return tmp
}

func (c CoinInquirer) GetBySymbol(symbol string) (*domains.Coin, error) {
	coin, err := c.client.GetBySymbol(c.ctx, &services.SymbolRequest{Symbol: symbol})
	if err != nil {
		return nil, err
	}
	createdAt, _ := time.Parse("2024-06-10", coin.CreatedAt)
	updatedAt, _ := time.Parse("2024-06-10", coin.UpdatedAt)
	return domains.NewCoin(coin.Id, coin.Name, coin.Symbol, coin.Icon, createdAt, updatedAt), nil
}

func NewCoinQuery(url string, ctx context.Context) query.CoinQuery {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error during connect : %v", err)
	}
	coinServiceClient := services.NewCoinServiceClient(conn)

	return &CoinInquirer{client: coinServiceClient, ctx: ctx}
}
