package user

import (
	"fmt"
	"go_grpc_realtime/lib/core/database"
	"go_grpc_realtime/lib/core/grpcgen"
	"go_grpc_realtime/lib/core/jwtmanager"
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

func (repo *repository) signUp(req *grpcgen.SignUpRequest) (*grpcgen.SignUpResponse, error) {
	if valErr := repo.Validation.ValidateEditUserRequest(req); valErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			valErr.Error(),
		)
	}

	hashPass, hashErr := utils.GenerateHashPassword(strings.TrimSpace(req.GetPassword()))
	if hashErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			hashErr.Error(),
		)
	}

	usr := UserTbl{
		FullName: strings.TrimSpace(req.GetUser().GetFullName()),
		Email:    strings.TrimSpace(req.GetUser().GetEmail()),
		Password: hashPass,
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

	return &grpcgen.SignUpResponse{
		User: &grpcgen.User{
			Id:       fmt.Sprint(usr.ID),
			FullName: usr.FullName,
			Email:    usr.Email,
		},
		JwtToken: jwtTkn,
	}, nil
}

func (repo *repository) loginUp(req *grpcgen.LoginRequest) (*grpcgen.SignUpResponse, error) {
	var usr UserTbl
	if err := database.DB.Where("email = ?", req.GetEmail()).First(&usr).Error; err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"user not found",
		)
	}

	isPasswordCor := utils.CheckPasswordCorHash(req.GetPassword(), usr.Password)
	if !isPasswordCor {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"invalid email or password",
		)
	}

	jwtTkn, jwtErr := jwtmanager.CreateToken(usr.ID)
	if jwtErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			jwtErr.Error(),
		)
	}

	return &grpcgen.SignUpResponse{
		User: &grpcgen.User{
			Id:       fmt.Sprint(usr.ID),
			FullName: usr.FullName,
			Email:    usr.Email,
		},
		JwtToken: jwtTkn,
	}, nil
}

func (repo *repository) getUsers(req *grpcgen.GetUsersRequest) (*grpcgen.GetUsersResponse, error) {
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

	var returnUsers []*grpcgen.User

	for _, user := range users {
		returnUsers = append(returnUsers, &grpcgen.User{
			Id:       fmt.Sprint(user.ID),
			FullName: user.FullName,
			Email:    user.Email,
		})
	}

	return &grpcgen.GetUsersResponse{
		Users: returnUsers,
	}, nil
}

func (repo *repository) updateUser(req *grpcgen.UpdateUserRequest, userId uint) (*grpcgen.User, error) {
	var usr UserTbl
	if err := database.DB.Where("id = ?", userId).First(&usr).Error; err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"User not found",
		)
	}

	updateBody := map[string]interface{}{}

	if req.GetFullName() != "" {
		if err := repo.Validation.IsStringValid(req.GetFullName()); err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"full name is empty",
			)
		}

		updateBody["full_name"] = strings.TrimSpace(req.GetFullName())
	}

	if req.GetEmail() != "" {
		if err := repo.Validation.IsEmail(req.GetEmail()); err != nil {
			return nil, status.Errorf(
				codes.Internal,
				err.Error(),
			)
		}

		updateBody["email"] = strings.TrimSpace(req.GetEmail())
	}

	if req.GetPassword() != "" {
		if err := repo.Validation.IsPasswordValid(req.GetPassword()); err != nil {
			return nil, status.Errorf(
				codes.Internal,
				err.Error(),
			)
		}

		hashPass, hashErr := utils.GenerateHashPassword(strings.TrimSpace(req.GetPassword()))
		if hashErr != nil {
			return nil, status.Errorf(
				codes.Internal,
				hashErr.Error(),
			)
		}

		updateBody["password"] = hashPass
	}

	if err := database.DB.Model(&usr).Updates(updateBody).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return &grpcgen.User{
		Id:       fmt.Sprint(usr.ID),
		FullName: usr.FullName,
		Email:    usr.Email,
	}, nil
}
