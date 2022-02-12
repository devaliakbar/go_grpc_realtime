package utils

import (
	"errors"
	"go_grpc_realtime/lib/core/generated/userpb"
	"net/mail"
	"strings"
)

type Validation struct{}

func (val *Validation) ValidateEditUserRequest(req *userpb.SignUpRequest) error {
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
