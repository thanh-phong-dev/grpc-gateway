package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
)

func PopulateToken() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = attachRequestID(ctx)
		isUserServiceMethodAuth := userServiceMethods()
		if isUserServiceMethodAuth[method] {
			return invoker(attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func userServiceMethods() map[string]bool {
	const userServicePath = "/user.UserService/"

	return map[string]bool{
		userServicePath + "GetUserInfo":   true,
		userServicePath + "GetUserDetail": false,
		userServicePath + "Register":      false,
	}
}

func attachToken(ctx context.Context) context.Context {
	userServiceToken := os.Getenv("USER_SERVICE_TOKEN")
	return metadata.AppendToOutgoingContext(ctx, "authorization", userServiceToken)
}

func attachRequestID(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	var (
		requestIdMD = md["requestid"]
		reqID       string
	)
	if len(requestIdMD) != 0 {
		reqID = requestIdMD[0]
	}

	return metadata.AppendToOutgoingContext(ctx, "requestid", reqID)
}
