package grpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"my-stocks/api-gateway/query"
	"my-stocks/common/grpc/services"
)

type AuthInquirer struct {
	client services.AuthServiceClient
	ctx    context.Context
}

func NewAuthQuery(url string, ctx context.Context) query.AuthQuery {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error during connect auth : %v", err)
	}
	authClient := services.NewAuthServiceClient(conn)

	return &AuthInquirer{client: authClient, ctx: ctx}
}

func (a AuthInquirer) VerifyToken(token string) (string, error) {
	response, err := a.client.VerifyToken(a.ctx, &services.VerifyRequest{Token: token})
	if err != nil {
		return "", err
	}

	return response.UserId, nil
}
