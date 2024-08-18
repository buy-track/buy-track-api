package grpc

import (
	"context"
	"my-stocks/auth/app"
	"my-stocks/common/grpc/services"
)

type authServer struct {
	authService app.AuthService
	services.UnimplementedAuthServiceServer
}

func newAuthServer(authService app.AuthService) *authServer {
	return &authServer{authService: authService}
}

func (s *authServer) VerifyToken(ctx context.Context, request *services.VerifyRequest) (*services.VerifyResponse, error) {
	found, err := s.authService.VerifyAccessToken(request.Token)
	if err != nil {
		return nil, err
	}
	return &services.VerifyResponse{UserId: found.UserId}, nil
}
