package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"realworld.tayler.io/internal/validator"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)

type CreateCommentDTO struct {
	Comment struct {
		Body *string `json:"body"`
	} `json:"comment"`
}

type Comment struct {
	CommentId int       `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    *Profile  `json:"author"`
}

func (comment CreateCommentDTO) Validate(v *validator.Validator) {
	v.Check(comment.Comment.Body != nil && *comment.Comment.Body != "", "body", "must not be empty")
}

type CommentRepository struct {
	DB             *sql.DB
	TimeoutSeconds int
}

func (repo *CommentRepository) CreateComment(articleId, currentUserId int, body string) (*Comment, error) {

	query := `INSERT INTO Comment (UserId, ArticleId, Body, CreatedAt, UpdatedAt) VALUES ($1, $2, $3, $4, $5) RETURNING CommentId`

	now := time.Now().UTC().Format(time.RFC3339Nano)

	args := []any{
		currentUserId,
		articleId,
		body,
		now,
		now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	var commentId int
	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&commentId)
	if err != nil {
		return nil, err
	}

	comment, err := repo.GetCommentById(commentId, currentUserId)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (repo *CommentRepository) GetCommentById(commentId, currentUserId int) (*Comment, error) {
	query := `SELECT 
	c.CommentId,
	c.Body,
	c.CreatedAt,
	c.UpdatedAt,
	EXISTS (SELECT 1 FROM Follower WHERE UserId = $1 AND FollowUserId = c.UserId) AS Following,
	u.Username,
	u.Bio,
	u.Image
  	FROM Comment c
  	JOIN User u ON c.UserId = u.UserId 
  	WHERE c.CommentId = $2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	var comment Comment
	var author Profile
	var createdAt string
	var updatedAt string

	err := repo.DB.QueryRowContext(ctx, query, currentUserId, commentId).Scan(
		&comment.CommentId,
		&comment.Body,
		&createdAt,
		&updatedAt,
		&author.Following,
		&author.Username,
		&author.Bio,
		&author.Image,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrCommentNotFound
		default:
			return nil, fmt.Errorf("an error occured while retrieving a comment: %w", err)
		}
	}

	comment.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing created at date: %w", err)
	}

	comment.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing updated at date: %w", err)
	}

	comment.Author = &author

	return &comment, nil
}
