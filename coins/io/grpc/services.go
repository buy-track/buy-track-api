package grpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"my-stocks/coins/app"
	"my-stocks/common/grpc/services"
)

type coinServer struct {
	coinService app.CoinService
	services.UnimplementedCoinServiceServer
}

func (c coinServer) Paginate(ctx context.Context, request *services.CoinQueryRequest) (*services.CoinPaginate, error) {
	coins := c.coinService.Paginate(int(request.Limit), int(request.Page))
	tmp := make([]*services.Coin, len(coins.Data))
	for i, coin := range coins.Data {
		tmp[i] = &services.Coin{
			CreatedAt: coin.CreatedAt.String(),
			UpdatedAt: coin.UpdatedAt.String(),
			Type:      int32(coin.Type),
			Price:     float32(coin.Price),
			Symbol:    coin.Symbol,
			Name:      coin.Name,
			Icon:      coin.Icon,
			Id:        coin.ID,
		}
	}
	return &services.CoinPaginate{
		Coins: &services.CoinList{Coins: tmp},
		Limit: request.Limit,
		Page:  request.Page,
	}, nil
}

func (c coinServer) GetByIds(ctx context.Context, request *services.IdsRequest) (*services.CoinList, error) {
	coins := c.coinService.GetByIds(request.Id)
	tmp := make([]*services.Coin, len(coins))
	for i, coin := range coins {
		tmp[i] = &services.Coin{
			CreatedAt: coin.CreatedAt.String(),
			UpdatedAt: coin.UpdatedAt.String(),
			Type:      int32(coin.Type),
			Price:     float32(coin.Price),
			Symbol:    coin.Symbol,
			Name:      coin.Name,
			Icon:      coin.Icon,
			Id:        coin.ID,
		}
	}
	return &services.CoinList{Coins: tmp}, nil
}

func (c coinServer) GetById(ctx context.Context, request *services.IdRequest) (*services.Coin, error) {
	coin, err := c.coinService.GetById(request.Id)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.NotFound, "Coin Not Found")
	}
	return &services.Coin{
		CreatedAt: coin.CreatedAt.String(),
		UpdatedAt: coin.UpdatedAt.String(),
		Type:      int32(coin.Type),
		Price:     float32(coin.Price),
		Symbol:    coin.Symbol,
		Name:      coin.Name,
		Icon:      coin.Icon,
		Id:        coin.ID,
	}, nil
}

func (c coinServer) GetBySymbol(ctx context.Context, request *services.SymbolRequest) (*services.Coin, error) {
	coin, err := c.coinService.GetBySymbol(request.Symbol)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Coin Not found")
	}
	return &services.Coin{
		CreatedAt: coin.CreatedAt.String(),
		UpdatedAt: coin.UpdatedAt.String(),
		Type:      int32(coin.Type),
		Price:     float32(coin.Price),
		Symbol:    coin.Symbol,
		Name:      coin.Name,
		Icon:      coin.Icon,
		Id:        coin.ID,
	}, nil
}

func newCoinServer(coinService app.CoinService) *coinServer {
	return &coinServer{coinService: coinService}
}
