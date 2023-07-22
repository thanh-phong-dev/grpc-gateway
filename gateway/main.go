package main

import (
	"gateway/internal/controller/authen"
	"gateway/internal/controller/user"
	"gateway/middleware"
	"gateway/pb/authpb"
	"gateway/pb/userpb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
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
	r := gin.Default()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging())

	userClientCreds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	if err != nil {
		log.Fatalf("[userClientCreds] get credentials got err: %v", err)
	}
	var userConn *grpc.ClientConn
	userConn, err = grpc.Dial(userServicePort,
		grpc.WithTransportCredentials(userClientCreds),
		grpc.WithUnaryInterceptor(middleware.PopulateToken()))
	if err != nil {
		log.Fatalf("[getUserServiceClient] did not connect: %s", err)
	}
	defer userConn.Close()

	authenClientCreds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	if err != nil {
		log.Fatalf("[authenClientCreds] get credentials got err: %v", err)
	}
	var authConn *grpc.ClientConn
	authConn, err = grpc.Dial(authServicePort,
		grpc.WithTransportCredentials(authenClientCreds),
		grpc.WithUnaryInterceptor(middleware.PopulateRequestID()),
	)
	if err != nil {
		log.Fatalf("[getAuthServiceClient] did not connect: %s", err)
	}
	defer authConn.Close()

	var (
		authService = authpb.NewAuthenticationServiceClient(authConn)
		userService = userpb.NewUserServiceClient(userConn)
	)

	authen.Route(r, authService, userService)
	user.Route(r, userService)

	r.Run(":8080")
}
