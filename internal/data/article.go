package data

import (
	"database/sql"
	"net/http"
	"time"

	"realworld.tayler.io/internal/validator"
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

type ArticleRepository struct {
	DB             *sql.DB
	TimeoutSeconds int
}

func (repo *ArticleRepository) CreateArticle(articleDto CreateArticleDTO) (*Article, error) {

	article := &Article{
		BodylessArticle: BodylessArticle{
			Title:       *articleDto.Article.Title,
			Description: *articleDto.Article.Description,
			TagList:     articleDto.Article.TagList,
		},
		Body: *articleDto.Article.Body,
	}

	return article, nil
}
