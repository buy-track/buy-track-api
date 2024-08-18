package main

import (
	"context"
	"github.com/golobby/container/v3"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"my-stocks/coins/app"
	"my-stocks/coins/cmd/ws_client"
	coinConfig "my-stocks/coins/config"
	"my-stocks/coins/io/grpc"
	"my-stocks/coins/persistance/postgres"
	redisRepository "my-stocks/coins/persistance/redis"
	"my-stocks/coins/persistance/repositories"
	"my-stocks/common/config"
)

var (
	cfg       *coinConfig.App
	conn      *gorm.DB
	redisDb   *redis.Client
	ctx       context.Context
	graceFull chan bool
	ctr       container.Container
)

func init() {
	ctx = context.Background()
	ctr = container.New()
	graceFull = make(chan bool, 1)

	var err error
	cfg, err = config.LoadConfig[coinConfig.App]()
	if err != nil {
		log.Fatalf("Error LoadConfig: %v", err.Error())
	}

	logLevel, _ := log.ParseLevel(cfg.LogLevel)
	log.SetLevel(logLevel)
}

func main() {

	connectDb()
	resolveServices()
	go grpc.Start(*cfg, ctr)

	var coinService app.CoinService
	_ = ctr.Resolve(&coinService)
	go ws_client.RunUpdateCoinPrice("wss://stream.binance.com:9443/ws", coinService)
	<-graceFull
}

func connectDb() {
	redisDb = redisRepository.GetConnection(cfg.RedisDB)
	conn = postgres.GetConnection(cfg.Database)
	postgres.Migrate(conn)
	resolveRepositories()

}

func resolveRepositories() {
	// Coins
	_ = ctr.Singleton(func() repositories.CoinProvider {
		return postgres.NewCoinRepository(conn)
	})

	_ = ctr.Singleton(func(provider repositories.CoinProvider) repositories.CoinReader {
		return provider
	})

	_ = ctr.Singleton(func(provider repositories.CoinProvider) repositories.CoinWriter {
		return provider
	})

	// CoinPrice
	_ = ctr.Singleton(func() repositories.CoinPriceProvider {
		return redisRepository.NewCoinPriceRepository(redisDb, ctx, cfg.RedisDB.Prefix)
	})

	_ = ctr.Singleton(func(provider repositories.CoinPriceProvider) repositories.CoinPriceReader {
		return provider
	})

	_ = ctr.Singleton(func(provider repositories.CoinPriceProvider) repositories.CoinPriceWriter {
		return provider
	})

}

func resolveServices() {
	// CoinService
	_ = ctr.Singleton(func(repository repositories.CoinProvider, priceRepo repositories.CoinPriceProvider) app.CoinService {
		return *app.NewCoinService(repository, repository, priceRepo, priceRepo)
	})
}
