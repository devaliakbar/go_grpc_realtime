package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func GetUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("-->Auth unary interceptor: ", info.FullMethod)

		return handler(ctx, req)
	}
}

func GetStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("-->Auth stream interceptor: ", info.FullMethod)

		return handler(srv, stream)
	}
}
