package user

import (
	"fmt"
	"go_grpc_realtime/lib/core/database"
	"go_grpc_realtime/lib/core/generated/userpb"
	"go_grpc_realtime/lib/core/jwtmanager"
	"go_grpc_realtime/lib/core/utils"
	"strconv"
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

func (repo *repository) signUp(req *userpb.SignUpRequest) (*userpb.SignUpResponse, error) {
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

	jwtTkn, jwtErr := jwtmanager.CreateToken(usr.ID)
	if jwtErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			jwtErr.Error(),
		)
	}

	return &userpb.SignUpResponse{
		User: &userpb.User{
			Id:       fmt.Sprint(usr.ID),
			FullName: usr.FullName,
			Email:    usr.Email,
		},
		JwtToken: jwtTkn,
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

func (repo *repository) updateUser(req *userpb.SignUpRequest) (*userpb.User, error) {
	userId, err := strconv.Atoi(req.GetUser().GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"User not found",
		)
	}

	var usr UserTbl
	if err := database.DB.Where("id = ?", uint(userId)).First(&usr).Error; err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"User not found",
		)
	}

	updateBody := map[string]interface{}{}

	if req.GetUser().FullName != "" {
		if err := repo.Validation.IsStringValid(req.GetUser().GetFullName()); err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"full name is empty",
			)
		}

		updateBody["full_name"] = strings.TrimSpace(req.GetUser().GetFullName())
	}

	if req.GetUser().GetEmail() != "" {
		if err := repo.Validation.IsEmail(req.GetUser().GetEmail()); err != nil {
			return nil, status.Errorf(
				codes.Internal,
				err.Error(),
			)
		}

		updateBody["email"] = strings.TrimSpace(req.GetUser().GetEmail())
	}

	if req.GetPassword() != "" {
		if err := repo.Validation.IsPasswordValid(req.GetPassword()); err != nil {
			return nil, status.Errorf(
				codes.Internal,
				err.Error(),
			)
		}

		updateBody["password"] = strings.TrimSpace(req.GetPassword())
	}

	if err := database.DB.Model(&usr).Updates(updateBody).Error; err != nil {
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
