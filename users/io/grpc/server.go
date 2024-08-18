package grpc

import (
	"github.com/golobby/container/v3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"my-stocks/common/grpc/services"
	"my-stocks/users/app"
	"my-stocks/users/config"
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
	var userService app.UserService
	_ = ctr.Resolve(&userService)
	userSrv := newUserServer(userService)
	services.RegisterUserServiceServer(grpcServer, userSrv)
}
