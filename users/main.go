package main

import (
	"context"
	"fmt"
	"github.com/golobby/container/v3"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"my-stocks/common/broker"
	"my-stocks/common/config"
	"my-stocks/users/app"
	userConfig "my-stocks/users/config"
	"my-stocks/users/io/event"
	"my-stocks/users/io/grpc"
	"my-stocks/users/persistance/postgres"
	"my-stocks/users/persistance/repositories"
)

var (
	cfg          *userConfig.App
	conn         *gorm.DB
	brkRedisConn *redis.Client
	eventDriven  *broker.EventDriven
	ctx          context.Context
	graceFull    chan bool
	ctr          container.Container
)

func init() {
	ctx = context.Background()
	ctr = container.New()
	graceFull = make(chan bool, 1)

	var err error
	cfg, err = config.LoadConfig[userConfig.App]()
	if err != nil {
		log.Fatalf("Error LoadConfig: %v", err.Error())
	}

	logLevel, _ := log.ParseLevel(cfg.LogLevel)
	log.SetLevel(logLevel)
}

func main() {

	connectDb()
	resolveServices()
	connectRedis()
	resolveBroker()
	go startEventDriven()
	go grpc.Start(*cfg, ctr)
	<-graceFull
}

func connectDb() {
	conn = postgres.GetConnection(cfg.Database)
	postgres.Migrate(conn)
	resolveRepositories()
}

func resolveRepositories() {
	// User
	_ = ctr.Singleton(func() repositories.UserProvider {
		return postgres.NewUserRepository(conn)
	})

	_ = ctr.Singleton(func(provider repositories.UserProvider) repositories.UserWriter {
		return provider
	})

	_ = ctr.Singleton(func(provider repositories.UserProvider) repositories.UserReader {
		return provider
	})

	// ProviderToken
	_ = ctr.Singleton(func() repositories.ProviderTokenProvider {
		return postgres.NewProviderTokenRepository(conn)
	})

	_ = ctr.Singleton(func(provider repositories.ProviderTokenProvider) repositories.ProviderTokenWriter {
		return provider
	})

	_ = ctr.Singleton(func(provider repositories.ProviderTokenProvider) repositories.ProviderTokenReader {
		return provider
	})

}

func resolveServices() {
	// UserService
	_ = ctr.Singleton(func(repository repositories.UserProvider) app.UserService {
		return *app.NewUserService(repository, repository)
	})

	// ProviderTokenService
	_ = ctr.Singleton(func(repository repositories.ProviderTokenProvider) app.ProviderTokenService {
		return *app.NewProviderTokenService(repository, repository)
	})
}

func connectRedis() {
	brkRedisConn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Broker.Host, cfg.Broker.Port),
		Password: cfg.Broker.Password,
		DB:       cfg.Broker.DB,
	})
}

func resolveBroker() {
	eventDriven = broker.New(broker.NewRedisBroker(brkRedisConn, ctx), ctx, graceFull)
}

func startEventDriven() {
	event.AddRoutes(eventDriven, ctr)
	event.Start()
}
