package grpc

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"my-stocks/api-gateway/query"
	"my-stocks/common/grpc/services"
	"my-stocks/domains"
	"time"
)

type UserInquirer struct {
	client services.UserServiceClient
	ctx    context.Context
}

func (u UserInquirer) EmailExists(email string) bool {
	exists, err := u.client.EmailExists(u.ctx, &services.EmailRequest{Email: email})
	if err != nil {
		log.Errorf("Error during check email : %v", err)
		return false
	}
	return exists.Exists
}

func (u UserInquirer) GetByEmail(email string) (*domains.User, error) {
	user, err := u.client.GetByEmail(u.ctx, &services.EmailRequest{Email: email})
	if err != nil {
		return nil, err
	}
	createdAt, _ := time.Parse("2024-06-10", user.CreatedAt)
	updatedAt, _ := time.Parse("2024-06-10", user.UpdatedAt)
	return &domains.User{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Email:     user.Email,
		Name:      user.Name,
		ID:        user.Id,
	}, nil
}

func (u UserInquirer) GetById(id string) (*domains.User, error) {
	user, err := u.client.GetById(u.ctx, &services.IdRequest{Id: id})
	if err != nil {
		return nil, err
	}
	createdAt, _ := time.Parse("2024-06-10", user.CreatedAt)
	updatedAt, _ := time.Parse("2024-06-10", user.UpdatedAt)
	return &domains.User{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Email:     user.Email,
		Name:      user.Name,
		ID:        user.Id,
	}, nil
}

func (u UserInquirer) CheckPassword(email string, password string) (*domains.User, error) {
	user, err := u.client.CheckPassword(u.ctx, &services.CheckPasswordRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	createdAt, _ := time.Parse("2024-06-10", user.CreatedAt)
	updatedAt, _ := time.Parse("2024-06-10", user.UpdatedAt)
	return &domains.User{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Email:     user.Email,
		Name:      user.Name,
		ID:        user.Id,
	}, nil
}

func NewUserQuery(url string, ctx context.Context) query.UserQuery {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error during connect : %v", err)
	}
	userServiceClient := services.NewUserServiceClient(conn)

	return &UserInquirer{client: userServiceClient, ctx: ctx}
}
