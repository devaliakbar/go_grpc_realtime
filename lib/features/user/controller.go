package user

import (
	"context"
	"go_grpc_realtime/lib/core/generated/userpb"
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

func (ctr *UserController) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	return ctr.repository.createUser(req)
}

func (ctr *UserController) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	return ctr.getUsers(req)
}

func (ctr *UserController) UpdateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	return ctr.repository.updateUser(req)
}
