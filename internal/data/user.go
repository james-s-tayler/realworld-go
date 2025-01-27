package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"realworld.tayler.io/internal/validator"
)

type User struct {
	UserId   int      `json:"-"`
	Email    string   `json:"email"`
	Token    string   `json:"token"`
	Username string   `json:"username"`
	Bio      string   `json:"bio"`
	Image    *string  `json:"image"`
	Password Password `json:"-"`
}

type Profile struct {
	Username  string  `json:"username"`
	Bio       string  `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"`
}

func (u *User) Validate(v *validator.Validator) {
	v.Check(v.Matches(u.Email, validator.EmailRX), "email", "must be a valid email address")
	v.Check(u.Username != "", "username", "must not be empty")
	v.Check(u.Bio != "", "bio", "must not be empty")

	if u.Password.Plaintext != nil {
		v.Check(len(*u.Password.Plaintext) >= 8, "password", "password must contain at least 8 characters")
	}

	if u.Password.hash == nil {
		// should never get here
		panic(fmt.Sprintf("missing password hash for user %v", u.Username))
	}
}

var (
	ErrDuplicateUsername  = errors.New("duplicate username")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type UserRepository struct {
	DB             *sql.DB
	TimeoutSeconds int
}

func (repo *UserRepository) RegisterUser(user *User) (*User, error) {

	query := `INSERT INTO USER (Email, Username, PasswordHash, Bio) VALUES($1,$2,$3,$4) RETURNING UserId`
	args := []any{user.Email, user.Username, user.Password.hash, user.Bio}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&user.UserId)
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

	query := `SELECT UserId, Username, Bio, Image, PasswordHash FROM User WHERE Email = $1`
	args := []any{email}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	user := &User{
		Email: email,
	}

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.UserId,
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

func (repo *UserRepository) GetUserById(userId int) (*User, error) {
	query := `SELECT Username, Email, Bio, Image, PasswordHash FROM User WHERE UserId = $1`
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	user := &User{
		UserId: userId,
	}

	err := repo.DB.QueryRowContext(ctx, query, userId).Scan(
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.Image,
		&user.Password.hash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("an unexpected error occurred when retrieving the user: %w", err)
		}
	}

	return user, nil
}

func (repo *UserRepository) GetUserByUsername(username string) (*User, error) {
	query := `SELECT UserId, Email, Bio, Image, PasswordHash FROM User WHERE Username = $1`
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	user := &User{
		Username: username,
	}

	err := repo.DB.QueryRowContext(ctx, query, username).Scan(
		&user.UserId,
		&user.Email,
		&user.Bio,
		&user.Image,
		&user.Password.hash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		default:
			return nil, fmt.Errorf("an unexpected error occurred when retrieving the user: %w", err)
		}
	}

	return user, nil
}

func (repo *UserRepository) UpdateUser(user *User) error {
	query := `UPDATE User SET (Username, Email, PasswordHash, Bio, Image) = ($1, $2, $3, $4, $5) WHERE UserId = $6`
	args := []any{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Bio,
		user.Image,
		user.UserId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == "UNIQUE constraint failed: User.Username":
			return ErrDuplicateUsername
		case err.Error() == "UNIQUE constraint failed: User.Email":
			return ErrDuplicateEmail
		default:
			return fmt.Errorf("error updating user: %w", err)
		}
	}
	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return fmt.Errorf("error updating user - no rows were updated: %w", err)
	}

	return nil
}
