package middleware

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
)

func Auth(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	err := authorize(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func authorize(ctx context.Context, method string) error {
	isAuthMethod := authMethods()
	if !isAuthMethod[method] {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("missing credentials")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	var (
		userServiceToken = os.Getenv("USER_SERVICE_TOKEN")
		accessToken      = values[0]
	)
	if accessToken != userServiceToken {
		return status.Error(codes.PermissionDenied, "no permission to access this RPC")
	}

	return nil
}

func authMethods() map[string]bool {
	const laptopServicePath = "/user.UserService/"

	return map[string]bool{
		laptopServicePath + "GetUserInfo": true,
		laptopServicePath + "Register":    false,
	}
}
