package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int      `json:"-"`
	Email    string   `json:"email"`
	Token    string   `json:"token"`
	Username string   `json:"username"`
	Bio      string   `json:"bio"`
	Image    *string  `json:"image"`
	Password Password `json:"-"`
}

var (
	ErrDuplicateUsername  = errors.New("duplicate username")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserRepository struct {
	DB             *sql.DB
	TimeoutSeconds int
}

func (repo *UserRepository) RegisterUser(user *User) (*User, error) {

	query := `INSERT INTO USER (Email, Username, PasswordHash, Bio) VALUES($1,$2,$3,$4) RETURNING Id`
	args := []any{user.Email, user.Username, user.Password.hash, user.Bio}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id)
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

func (repo *UserRepository) GetUserByCredentials(email string, password string) (*User, error) {

	query := `SELECT Username, Bio, Image, PasswordHash FROM User WHERE Email = $1`
	args := []any{email}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	user := &User{
		Email: email,
	}

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.Username,
		&user.Bio,
		&user.Image,
		&user.Password.hash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrInvalidCredentials
		default:
			return nil, fmt.Errorf("error when looking up user for credential check: %w", err)
		}
	}

	err = bcrypt.CompareHashAndPassword(user.Password.hash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return nil, ErrInvalidCredentials
		default:
			return nil, fmt.Errorf("error when attempting to compare password and password hash: %w", err)
		}
	}

	return user, nil
}
