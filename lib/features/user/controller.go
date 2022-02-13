package user

import (
	"context"
	"go_grpc_realtime/lib/core/grpcgen"
	"go_grpc_realtime/lib/core/jwtmanager"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct {
	grpcgen.UnimplementedUserServiceServer
	*repository
}

func InitAndGetUserServices() grpcgen.UserServiceServer {
	repo := &repository{}

	repo.migrateDb()

	return &UserController{
		repository: repo,
	}
}

func (ctr *UserController) SignUp(ctx context.Context, req *grpcgen.SignUpRequest) (*grpcgen.SignUpResponse, error) {
	return ctr.repository.signUp(req)
}

func (ctr *UserController) GetUsers(ctx context.Context, req *grpcgen.GetUsersRequest) (*grpcgen.GetUsersResponse, error) {
	return ctr.getUsers(req)
}

func (ctr *UserController) UpdateUser(ctx context.Context, req *grpcgen.UpdateUserRequest) (*grpcgen.User, error) {
	userId, ok := ctx.Value(jwtmanager.USER_ID_KEY).(uint)

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"User not found",
		)
	}

	return ctr.repository.updateUser(req, userId)
}

func (ctr *UserController) Login(ctx context.Context, req *grpcgen.LoginRequest) (*grpcgen.SignUpResponse, error) {
	return ctr.loginUp(req)
}
