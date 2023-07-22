package service

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"user/internal/model"
	"user/internal/repository"
	"user/pb/userpb"
)

type authentication struct {
	user     string
	password string
}

type UserGrpcServer struct {
	userpb.UserServiceServer
	userRepo repository.UserRepository
	auth     *authentication
}

func NewGrpcUserServer(userRepo repository.UserRepository) *UserGrpcServer {
	return &UserGrpcServer{
		userRepo: userRepo,
		auth: &authentication{
			user:     "phong",
			password: "phong123456",
		},
	}
}

func (authentication *authentication) Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("missing credentials")
	}

	var us, pw string
	if val, ok := md["user"]; ok {
		us = val[0]
	}
	if val, ok := md["password"]; ok {
		pw = val[0]
	}

	if us != authentication.user || pw != authentication.password {
		return errors.New("invalid credentials")
	}

	return nil
}

func (s *UserGrpcServer) Register(ctx context.Context, registerRequest *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	user := model.User{
		Name:     registerRequest.Name,
		Username: registerRequest.Username,
		Password: registerRequest.Password,
		Gender:   registerRequest.Gender,
		Email:    registerRequest.Email,
	}

	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &userpb.RegisterResponse{Status: "registered successfully"}, nil
}

func (s *UserGrpcServer) GetUserInfo(ctx context.Context, userInfoRequest *userpb.UserInfoRequest) (*userpb.UserInfoResponse, error) {
	userCondition := model.User{
		ID:       uint(userInfoRequest.Id),
		Username: userInfoRequest.Username,
	}

	user, err := s.userRepo.FindOne(userCondition)
	if err != nil {
		return nil, err
	}

	return &userpb.UserInfoResponse{
		Username: user.Username,
		Password: user.Password,
		Name:     user.Name,
		Gender:   user.Gender,
		Email:    user.Email,
	}, nil
}
