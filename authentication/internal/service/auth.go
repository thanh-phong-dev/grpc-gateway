package service

import (
	"authentication/pb/authpb"
	"authentication/pb/userpb"
	"context"
	"github.com/google/uuid"
)

type Server struct {
	authpb.AuthenticationServiceServer
	UserService userpb.UserServiceClient
}

func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	user, err := s.UserService.GetUserInfo(ctx, &userpb.UserInfoRequest{
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	if req.Password != user.Password {
		return &authpb.LoginResponse{
			ErrorMessage: "Password is incorrect",
		}, nil
	}

	return &authpb.LoginResponse{
		Token: user.Username + uuid.New().String(),
	}, nil
}
