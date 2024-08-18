package main

import (
	"context"
	"fmt"
	"github.com/golobby/container/v3"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"my-stocks/api-gateway/app"
	"my-stocks/api-gateway/command"
	"my-stocks/api-gateway/command/event"
	userConfig "my-stocks/api-gateway/config"
	"my-stocks/api-gateway/io/http"
	"my-stocks/api-gateway/query"
	"my-stocks/api-gateway/query/grpc"
	"my-stocks/common/broker"
	"my-stocks/common/config"
)

var (
	cfg              *userConfig.App
	ctx              context.Context
	authBrkRedisConn *redis.Client
	userBrkRedisConn *redis.Client
	graceFull        chan bool
	ctr              container.Container
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
	connectRedis()
	resolveCommands()
	resolveQueries()
	resolveServices()

	go http.Start(cfg, ctr)
	<-graceFull
}

func connectRedis() {
	authBrkRedisConn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.AuthBroker.Host, cfg.AuthBroker.Port),
		Password: cfg.AuthBroker.Password,
		DB:       cfg.AuthBroker.DB,
	})
	userBrkRedisConn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.UserBroker.Host, cfg.UserBroker.Port),
		Password: cfg.UserBroker.Password,
		DB:       cfg.UserBroker.DB,
	})
}

func resolveCommands() {
	// AuthCommand
	_ = ctr.Singleton(func() command.AuthCommander {
		return event.NewAuthCommand(broker.NewRedisBroker(authBrkRedisConn, ctx))
	})

	// UserCommand
	_ = ctr.Singleton(func() command.UserCommander {
		return event.NewUserCommand(broker.NewRedisBroker(userBrkRedisConn, ctx))
	})
}

func resolveQueries() {
	// AuthQuery
	_ = ctr.Singleton(func() query.AuthQuery {
		return grpc.NewAuthQuery(cfg.AuthServiceGrpcAddress, ctx)
	})

	// UserQuery
	_ = ctr.Singleton(func() query.UserQuery {
		return grpc.NewUserQuery(cfg.UserServiceGrpcAddress, ctx)
	})

	// CoinQuery
	_ = ctr.Singleton(func() query.CoinQuery {
		return grpc.NewCoinQuery(cfg.CoinServiceGrpcAddress, ctx)
	})

}

func resolveServices() {
	// AuthService
	_ = ctr.Singleton(func(command command.AuthCommander, query query.AuthQuery) app.AuthService {
		return *app.NewAuthService(command, query)
	})

	// UserService
	_ = ctr.Singleton(func(command command.UserCommander, query query.UserQuery) app.UserService {
		return *app.NewUserService(command, query)
	})

	// UserCoin
	_ = ctr.Singleton(func(query query.CoinQuery) app.CoinService {
		return *app.NewCoinService(nil, query)
	})
}
