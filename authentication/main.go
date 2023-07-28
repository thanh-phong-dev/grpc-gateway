package main

import (
	"authentication/internal/service"
	"authentication/middleware"
	"authentication/pb/authpb"
	"authentication/pb/userpb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

const (
	authServicePort = ":9000"
	userServicePort = ":9001"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load env file. Err: %s", err)
	}

	userClientCreds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	if err != nil {
		log.Fatalf("[userClientCreds] get credentials got err: %v", err)
	}

	var userServiceConn *grpc.ClientConn
	userServiceConn, err = grpc.Dial(userServicePort,
		grpc.WithTransportCredentials(userClientCreds),
		grpc.WithUnaryInterceptor(middleware.PopulateToken()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer userServiceConn.Close()

	authenServerCreds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatalf("[authenServerCreds] failed to get credential: %v", err)
	}
	var (
		userService       = userpb.NewUserServiceClient(userServiceConn)
		authServiceServer = service.Server{
			UserService: userService,
		}
		grpcServer = grpc.NewServer(
			grpc.Creds(authenServerCreds),
			grpc.UnaryInterceptor(middleware.Logging),
		)
	)

	authpb.RegisterAuthenticationServiceServer(grpcServer, &authServiceServer)

	lis, err := net.Listen("tcp", authServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
