package user

import (
	"fmt"
	"go_grpc_realtime/lib/core/database"
	"go_grpc_realtime/lib/core/generated/userpb"
	"go_grpc_realtime/lib/core/utils"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repository struct {
	*utils.Validation
}

func (*repository) migrateDb() {
	database.DB.AutoMigrate(&UserTbl{})
}

func (repo *repository) createUser(req *userpb.CreateUserRequest) (*userpb.User, error) {
	if valErr := repo.Validation.ValidateEditUserRequest(req); valErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			valErr.Error(),
		)
	}

	usr := UserTbl{
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

func (repo *repository) getUsers(req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	skip := int(req.GetSkip())

	take := 10
	if req.GetTake() != 0 {
		take = int(req.GetTake())
		if take > 100 {
			take = 100
		}
	}

	var opts []interface{}
	searchQry := strings.TrimSpace(req.GetSearch())
	if searchQry != "" {
		opts = append(opts, "full_name LIKE ? OR email LIKE ?")
		opts = append(opts, "%"+searchQry+"%")
		opts = append(opts, "%"+searchQry+"%")
	}

	var users []UserQuery

	database.DB.Model(&UserTbl{}).Order("full_name asc").Offset(skip).Limit(take).Find(&users, opts...)

	var returnUsers []*userpb.User

	for _, user := range users {
		returnUsers = append(returnUsers, &userpb.User{
			Id:       fmt.Sprint(user.ID),
			FullName: user.FullName,
			Email:    user.Email,
		})
	}

	return &userpb.GetUsersResponse{
		Users: returnUsers,
	}, nil
}
