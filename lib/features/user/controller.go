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

func (ctr *UserController) SignUp(ctx context.Context, req *userpb.SignUpRequest) (*userpb.SignUpResponse, error) {
	return ctr.repository.signUp(req)
}

func (ctr *UserController) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	return ctr.getUsers(req)
}

func (ctr *UserController) UpdateUser(ctx context.Context, req *userpb.SignUpRequest) (*userpb.User, error) {
	return ctr.repository.updateUser(req)
}
