package data

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ITokenService interface {
	CreateToken(user *User) (string, error)
}

type JwtTokenService struct {
	SecretKey []byte
}

func (t JwtTokenService) CreateToken(user *User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(t.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
