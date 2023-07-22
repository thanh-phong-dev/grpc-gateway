package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func Logging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	var (
		method      = info.FullMethod
		requestIdMD = md["requestid"]
		reqID       string
	)
	if len(requestIdMD) != 0 {
		reqID = requestIdMD[0]
	}

	resp, err := handler(ctx, req)

	log.Printf("[%v] RequestId: %s, Received resquest: %v, Received response: %v", method, reqID, req, resp)
	return resp, err
}
