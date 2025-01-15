package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id       int      `json:"-"`
	Email    string   `json:"email"`
	Token    string   `json:"token"`
	Username string   `json:"username"`
	Bio      string   `json:"bio"`
	Image    *string  `json:"image"`
	Password password `json:"-"`
}

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicateEmail    = errors.New("duplicate email")
)

type UserRepository struct {
	DB             *sql.DB
	TimeoutSeconds int
}

func (userRepo *UserRepository) RegisterUser(user *User) (*User, error) {

	query := `INSERT INTO USER (Email, Username, PasswordHash, Bio) VALUES($1,$2,$3,$4) RETURNING Id`
	args := []any{user.Email, user.Username, user.Password.hash, user.Bio}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(userRepo.TimeoutSeconds)*time.Second)
	defer cancel()

	err := userRepo.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id)
	if err != nil {
		switch {
		case err.Error() == "UNIQUE constraint failed: User.Username":
			return nil, ErrDuplicateUsername
		case err.Error() == "UNIQUE constraint failed: User.Email":
			return nil, ErrDuplicateEmail
		default:
			return nil, fmt.Errorf("error when registering a user: %w", err)
		}
	}

	return user, nil
}
