package data

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ITokenService interface {
	CreateToken(user *User) (string, error)
	VerifyToken(tokenString string) (*jwt.Token, error)
}

type JwtTokenService struct {
	SecretKey []byte
}

type CustomClaims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (t JwtTokenService) CreateToken(user *User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":  user.Id,
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(t.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t JwtTokenService) VerifyToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return t.SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("an unknown error occured while attempting to parse the token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
