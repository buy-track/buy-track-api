package grpc

import (
	"github.com/golobby/container/v3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"my-stocks/coins/app"
	"my-stocks/coins/config"
	"my-stocks/common/grpc/services"
	"net"
)

func Start(cfg config.App, ctr container.Container) {
	lis, err := net.Listen("tcp", cfg.GrpcServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	initServers(grpcServer, ctr)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}

func initServers(grpcServer *grpc.Server, ctr container.Container) {
	var coinService app.CoinService
	_ = ctr.Resolve(&coinService)

	coinSrv := newCoinServer(coinService)

	services.RegisterCoinServiceServer(grpcServer, coinSrv)
}
