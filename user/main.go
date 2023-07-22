package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"user/database"
	"user/internal/repository"
	"user/internal/service"
	"user/middleware"
	"user/pb/userpb"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load env file. Err: %s", err)
	}

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userServerCreds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatalf("[userServerCreds] failed to get credential: %v", err)
	}

	var (
		db                = database.Init()
		userRepo          = repository.NewUserRepository(db)
		userServiceServer = service.NewGrpcUserServer(userRepo)
		grpcServer        = grpc.NewServer(
			grpc.Creds(userServerCreds),
			grpc.ChainUnaryInterceptor(middleware.Auth, middleware.Logging),
		)
	)

	userpb.RegisterUserServiceServer(grpcServer, userServiceServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
