package user

import (
	"fmt"
	"go_grpc_realtime/lib/core/database"
	"go_grpc_realtime/lib/core/generated/userpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repository struct{}

func (*repository) migrateDb() {
	database.DB.AutoMigrate(&User{})
}

func (*repository) createUser(req *userpb.CreateUserRequest) (*userpb.User, error) {

	usr := User{
		FullName: req.GetUser().GetFullName(),
		Email:    req.GetUser().GetEmail(),
		Password: req.GetPassword(),
	}
	if err := database.DB.Create(&usr).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return &userpb.User{
		Id:       fmt.Sprint(usr.ID),
		FullName: req.GetUser().GetFullName(),
		Email:    req.GetUser().GetEmail(),
	}, nil
}
