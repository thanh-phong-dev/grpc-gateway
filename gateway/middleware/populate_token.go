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
		isUserServiceMethodAuth := userServiceMethods()
		ctx = attachRequestID(ctx)
		if isUserServiceMethodAuth[method] {
			return invoker(attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func PopulateRequestID() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(attachRequestID(ctx), method, req, reply, cc, opts...)
	}
}

func userServiceMethods() map[string]bool {
	const userServicePath = "/user.UserService/"

	return map[string]bool{
		userServicePath + "GetUserInfo": true,
		userServicePath + "Register":    false,
	}
}

func attachToken(ctx context.Context) context.Context {
	userServiceToken := os.Getenv("USER_SERVICE_TOKEN")
	return metadata.AppendToOutgoingContext(ctx, "authorization", userServiceToken)
}

func attachRequestID(ctx context.Context) context.Context {
	requestId, ok := ctx.Value("requestId").(string)
	if !ok {
		requestId = ""
	}
	return metadata.AppendToOutgoingContext(ctx, "requestid", requestId)
}
