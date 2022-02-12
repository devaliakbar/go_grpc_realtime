package user

import (
	"fmt"
	"go_grpc_realtime/lib/core/database"
	"go_grpc_realtime/lib/core/generated/userpb"
	"go_grpc_realtime/lib/core/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repository struct {
	*utils.Validation
}

func (*repository) migrateDb() {
	database.DB.AutoMigrate(&User{})
}

func (repo *repository) createUser(req *userpb.CreateUserRequest) (*userpb.User, error) {
	if valErr := repo.Validation.ValidateEditUserRequest(req); valErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			valErr.Error(),
		)
	}

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
		FullName: usr.FullName,
		Email:    usr.Email,
	}, nil
}
