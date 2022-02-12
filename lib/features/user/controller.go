package user

import (
	"context"
	"go_grpc_realtime/lib/core/generated/userpb"
	"log"
)

type UserController struct {
	userpb.UnimplementedUserServiceServer
	*repository
}

func InitAndGetUserServices() userpb.UserServiceServer {
	repo := &repository{}

	repo.migrateDb()

	return &UserController{
		repository: repo,
	}
}

func (controller *UserController) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	log.Printf("--->CreateUser: %v", req.GetUser())

	return controller.repository.createUser(req)
}
