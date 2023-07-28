package main

import (
	"context"
	"gateway/middleware"
	"gateway/pb/authpb"
	"gateway/pb/userpb"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load env file. Err: %s", err)
	}
	userClientCreds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	if err != nil {
		log.Fatalf("[userClientCreds] get credentials got err: %v", err)
	}

	authenClientCreds, err := credentials.NewClientTLSFromFile("cert/server.crt", "localhost")
	if err != nil {
		log.Fatalf("[authenClientCreds] get credentials got err: %v", err)
	}

	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			requestID, ok := ctx.Value("requestID").(string)
			if !ok {
				requestID = ""
			}
			md := metadata.Pairs("authorization", header, "requestid", requestID)
			return md
		}))

	err = authpb.RegisterAuthenticationServiceHandlerFromEndpoint(ctx, mux,
		"localhost:9000",
		[]grpc.DialOption{
			grpc.WithTransportCredentials(authenClientCreds),
			grpc.WithUnaryInterceptor(middleware.PopulateToken()),
		})
	if err != nil {
		log.Fatalf("[RegisterAuthenticationServiceHandlerFromEndpoint] got err: %v", err)
		return
	}

	err = userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux,
		"localhost:9001",
		[]grpc.DialOption{
			grpc.WithTransportCredentials(userClientCreds),
			grpc.WithUnaryInterceptor(middleware.PopulateToken()),
		},
	)
	if err != nil {
		log.Fatalf("[RegisterUserServiceHandlerFromEndpoint] got err: %v", err)
		return
	}

	// Creating a normal HTTP server
	server := gin.New()
	server.Use(middleware.RequestID())
	server.Use(middleware.Logging())
	server.Group("*{grpc_gateway}").Any("", gin.WrapH(mux))

	// start server
	err = server.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
