package utils

import (
	"errors"
	"go_grpc_realtime/lib/core/grpcgen"
	"net/mail"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validation struct{}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////***User Service Validation***////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (val *Validation) ValidateEditUserRequest(req *grpcgen.SignUpRequest) error {
	if err := val.IsStringValid(req.GetUser().GetFullName()); err != nil {
		return errors.New("full name is empty")
	}

	if err := val.IsEmail(req.GetUser().GetEmail()); err != nil {
		return err
	}

	if err := val.IsPasswordValid(req.GetPassword()); err != nil {
		return err
	}

	return nil
}

func (*Validation) IsEmail(email string) error {
	_, err := mail.ParseAddress(email)

	if err != nil {
		return errors.New("invalid mail address")
	}
	return nil
}

func (*Validation) IsStringValid(strg string) error {
	if strings.TrimSpace(strg) == "" {
		return errors.New("string is empty")
	}
	return nil
}

func (*Validation) IsPasswordValid(strg string) error {
	pasLen := len(strings.TrimSpace(strg))
	if pasLen < 6 || pasLen > 50 {
		return errors.New("password length must be between 6 and 50")
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////***Message Service Validation***////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func RemoveDuplicateUsers(usersStr []string) ([]uint, error) {
	keys := make(map[uint]bool)
	list := []uint{}

	for _, entry := range usersStr {
		userId, err := strconv.Atoi(entry)
		if err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"invalid user id",
			)
		}

		if value := keys[uint(userId)]; !value {
			keys[uint(userId)] = true
			list = append(list, uint(userId))
		}
	}
	return list, nil
}
