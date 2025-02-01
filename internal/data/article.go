package data

import (
	"net/http"
	"time"
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
