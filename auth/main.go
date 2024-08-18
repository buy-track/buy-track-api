package main

import (
	"context"
	"fmt"
	"github.com/golobby/container/v3"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"my-stocks/auth/app"
	authConfig "my-stocks/auth/config"
	"my-stocks/auth/io/event"
	"my-stocks/auth/io/grpc"
	"my-stocks/auth/persistance/postgres"
	"my-stocks/auth/persistance/repositories"
	"my-stocks/common/broker"
	"my-stocks/common/config"
)

var (
	cfg          *authConfig.App
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
	cfg, err = config.LoadConfig[authConfig.App]()
	fmt.Println(cfg)
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
	// AccessToken
	_ = ctr.Singleton(func() repositories.AccessTokenProvider {
		return postgres.NewAccessTokenRepository(conn)
	})

	_ = ctr.Singleton(func(provider repositories.AccessTokenProvider) repositories.AccessTokenReader {
		return provider
	})

	_ = ctr.Singleton(func(provider repositories.AccessTokenProvider) repositories.AccessTokenWriter {
		return provider
	})
}

func resolveServices() {
	// AuthService
	_ = ctr.Singleton(func(repository repositories.AccessTokenProvider) app.AuthService {
		return *app.NewAuthService(repository, repository)
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
