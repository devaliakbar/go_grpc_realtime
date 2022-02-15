package interceptors

import (
	"context"
	"go_grpc_realtime/lib/core/jwtmanager"
	"go_grpc_realtime/lib/core/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//All routes which auth is not required
var authNotRequiredRoutes = []string{
	"/user.UserService/SignUp",
	"/user.UserService/Login",
}

func GetUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		userId, err := checkAuthorization(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		if userId != 0 {
			ctx = context.WithValue(ctx, jwtmanager.USER_ID_KEY, userId)
		}
		return handler(ctx, req)
	}
}

func checkAuthorization(ctx context.Context, method string) (uint, error) {
	///If auth not required
	if utils.CheckStringExist(authNotRequiredRoutes, method) {
		return 0, nil
	}

	///If auth required
	return GetUserIdFromHeader(ctx)
}

func GetUserIdFromHeader(ctx context.Context) (uint, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return 0, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	userId, error := jwtmanager.IsTokenValid(accessToken)
	if error != nil {
		return 0, status.Error(codes.Unauthenticated, "invalid authorization token")
	}

	return userId, nil
}
