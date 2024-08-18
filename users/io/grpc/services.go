package grpc

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"my-stocks/common/grpc/services"
	"my-stocks/users/app"
)

type userServer struct {
	userService app.UserService
	services.UnimplementedUserServiceServer
}

func newUserServer(userService app.UserService) *userServer {
	return &userServer{userService: userService}
}

func (u *userServer) EmailExists(ctx context.Context, request *services.EmailRequest) (*services.EmailExistsResponse, error) {
	exists := u.userService.EmailExists(request.Email)
	return &services.EmailExistsResponse{Exists: exists}, nil
}

func (u *userServer) GetByEmail(ctx context.Context, request *services.EmailRequest) (*services.UserResponse, error) {
	user, err := u.userService.FindByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	return &services.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (u *userServer) GetById(ctx context.Context, request *services.IdRequest) (*services.UserResponse, error) {
	user, err := u.userService.FindById(request.Id)
	if err != nil {
		return nil, err
	}
	return &services.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (u *userServer) CheckPassword(ctx context.Context, request *services.CheckPasswordRequest) (*services.UserResponse, error) {
	user, err := u.userService.FindByEmail(request.Email)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "password or email is invalid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "password or email is invalid")
	}

	return &services.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
