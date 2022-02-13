package user

import (
	"context"
	"go_grpc_realtime/lib/core/grpc_generated/userpb"
	"go_grpc_realtime/lib/core/jwtmanager"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (ctr *UserController) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error) {
	userId, ok := ctx.Value(jwtmanager.USER_ID_KEY).(uint)

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"User not found",
		)
	}

	return ctr.repository.updateUser(req, userId)
}

func (ctr *UserController) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.SignUpResponse, error) {
	return ctr.loginUp(req)
}
