package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"realworld.tayler.io/internal/validator"
)

var (
	ErrArticleNotFound = errors.New("article not found")
)

type ArticleFilters struct {
	Tag       *string
	Author    *string
	Favorited *string
}

type BodylessArticle struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         *Profile  `json:"author"`
}

type Article struct {
	BodylessArticle
	Body string `json:"body"`
}

type CreateArticleDTO struct {
	Article struct {
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		Body        *string  `json:"body"`
		TagList     []string `json:"tagList"`
	} `json:"article"`
}

func (f *ArticleFilters) ParseFilters(r *http.Request) {
	if r.URL.Query().Has("tag") {
		value := r.URL.Query().Get("tag")
		f.Tag = &value
	}
	if r.URL.Query().Has("author") {
		value := r.URL.Query().Get("author")
		f.Author = &value
	}
	if r.URL.Query().Has("favorited") {
		value := r.URL.Query().Get("favorited")
		f.Favorited = &value
	}
}

func (article CreateArticleDTO) Validate(v *validator.Validator) {
	v.Check(article.Article.Body != nil && *article.Article.Body != "", "body", "must not be empty")
	v.Check(article.Article.Description != nil && *article.Article.Description != "", "description", "must not be empty")
	v.Check(article.Article.Title != nil && *article.Article.Title != "", "title", "must not be empty")
}

func (article CreateArticleDTO) GetSlug() string {
	return strings.ToLower(strings.ReplaceAll(*article.Article.Title, " ", "-"))
}

type ArticleRepository struct {
	DB             *sql.DB
	TimeoutSeconds int
}

func (repo *ArticleRepository) CreateArticle(articleDto CreateArticleDTO, userId int) (*Article, error) {

	query := `INSERT INTO Article 
				(UserId, Slug, Title, Description, Body, CreatedAt, UpdatedAt) 
				VALUES ($1, $2, $3, $4, $5, $6, $7) 
				RETURNING ArticleId`

	now := time.Now().UTC().Format(time.RFC3339Nano)

	args := []any{
		userId,
		articleDto.GetSlug(),
		*articleDto.Article.Title,
		*articleDto.Article.Description,
		*articleDto.Article.Body,
		now,
		now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	// start a transaction

	var articleId int
	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&articleId)
	if err != nil {
		return nil, fmt.Errorf("an error occurred when saving article: %w", err)
	}

	article, err := repo.GetArticleBySlug(articleDto.GetSlug())
	if err != nil {
		return nil, fmt.Errorf("an error occurred when looking up article by slug after saving: %w", err)
	}

	//end transaction

	return article, nil
}

func (repo *ArticleRepository) GetArticleBySlug(slug string) (*Article, error) {

	// need to save and load tags

	query := `SELECT 
				a.Title,
				a.Slug,
				a.Description,
				a.Body,
				a.CreatedAt,
				a.UpdatedAt,
				(SELECT EXISTS(SELECT 1 FROM ArticleFavorite WHERE ArticleId=a.ArticleId AND UserId=$1)) AS Favorited,
				(SELECT COUNT(*) FROM ArticleFavorite WHERE ArticleId=a.ArticleId) AS FavoritesCount,
				COALESCE((SELECT GROUP_CONCAT(t.Tag, ',') 
				        FROM Tag t 
                        JOIN ArticleTag at ON at.TagId = t.TagId 
                        WHERE at.ArticleId = a.ArticleId), '') AS Tags,
				u.Username,
				u.Bio,
				u.Image
			  FROM Article a
			  JOIN User u ON a.UserId = u.UserId 
			  WHERE a.Slug = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	var article Article
	var author Profile
	var rawTags string
	var createdAt string
	var updatedAt string

	err := repo.DB.QueryRowContext(ctx, query, slug).Scan(
		&article.Title,
		&article.Slug,
		&article.Description,
		&article.Body,
		&createdAt,
		&updatedAt,
		&article.Favorited,
		&article.FavoritesCount,
		&rawTags,
		&author.Username,
		&author.Bio,
		&author.Image,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrArticleNotFound
		default:
			return nil, fmt.Errorf("error looking up article by slug: %w", err)
		}
	}

	article.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing created at date: %w", err)
	}

	article.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing updated at date: %w", err)
	}

	article.Author = &author

	if rawTags == "" {
		article.TagList = make([]string, 0)
	} else {
		article.TagList = strings.Split(rawTags, ",")
	}

	return &article, nil
}
