package user

import (
	"context"
	"go_grpc_realtime/lib/core/generated/userpb"
	"log"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
}

func (*UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	log.Printf("--->CreateUser: %v", req.GetUser())

	return &userpb.User{
		FullName: req.GetUser().GetFullName(),
		Email:    req.GetUser().GetEmail(),
	}, nil
}
