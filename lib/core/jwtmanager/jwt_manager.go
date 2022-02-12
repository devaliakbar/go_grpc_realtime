package jwtmanager

import (
	"errors"
	"fmt"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
)

type userIdKeyType string

const USER_ID_KEY userIdKeyType = "userId"

const JWT_SECRET = "123456"

func CreateToken(user_id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	///claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET))
}

func IsTokenValid(token string) (uint, error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return 0, err
	}

	return extractUserIDfronToken(tokenObj)
}

func extractUserIDfronToken(token *jwt.Token) (uint, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, errors.New("can't find user associated with token")
}
