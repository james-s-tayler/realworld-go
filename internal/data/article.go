package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"realworld.tayler.io/internal/validator"
)

var (
	ErrArticleNotFound = errors.New("article not found")
	ErrDuplicateSlug   = errors.New("duplicate slug")
)

type ArticleFilters struct {
	Tag       *string
	Author    *string
	Favorited *string
	PaginationFilters
}

type PaginationFilters struct {
	Limit  int
	Offset int
}

type BodylessArticle struct {
	ArticleId      int       `json:"-"`
	UserId         int       `json:"-"`
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

type UpdateArticleDTO struct {
	Article struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Body        *string `json:"body"`
	} `json:"article"`
}

func (f *ArticleFilters) ParseFilters(v *validator.Validator, r *http.Request) {
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
	f.PaginationFilters = PaginationFilters{}
	f.PaginationFilters.ParseFilters(v, r)
}

func (f *PaginationFilters) ParseFilters(v *validator.Validator, r *http.Request) {
	if r.URL.Query().Has("limit") {
		value := r.URL.Query().Get("limit")
		limit, err := strconv.Atoi(value)
		if err != nil {
			v.AddError("limit", "must be an integer")
		} else {
			f.Limit = limit
		}
	} else {
		f.Limit = 20
	}

	if r.URL.Query().Has("offset") {
		value := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(value)
		if err != nil {
			v.AddError("offset", "must be an integer")
		} else {
			f.Offset = offset
		}
	} else {
		f.Offset = 0
	}

	v.Check(f.Limit > 0, "limit", "must be a positive integer")
	v.Check(f.Offset >= 0, "offset", "must be greater than or equal to zero")
}

func (article CreateArticleDTO) Validate(v *validator.Validator) {
	v.Check(article.Article.Body != nil && *article.Article.Body != "", "body", "must not be empty")
	v.Check(article.Article.Description != nil && *article.Article.Description != "", "description", "must not be empty")
	v.Check(article.Article.Title != nil && *article.Article.Title != "", "title", "must not be empty")
	for _, tag := range article.Article.TagList {
		v.Check(tag != "", "tagList", "tag must not be blank")
		v.Check(!strings.Contains(tag, ","), "tagList", "must not contain ','")
	}
}

func (article UpdateArticleDTO) Validate(v *validator.Validator) {

	v.Check(article.Article.Body != nil ||
		article.Article.Description != nil ||
		article.Article.Title != nil, "article", "must provide at least one of body, description or title")

	if article.Article.Body != nil {
		v.Check(*article.Article.Body != "", "body", "must not be blank")
	}
	if article.Article.Description != nil {
		v.Check(*article.Article.Description != "", "description", "must not be blank")
	}
	if article.Article.Title != nil {
		v.Check(*article.Article.Title != "", "title", "must not be blank")
	}
}

func (article CreateArticleDTO) GetSlug() string {
	return strings.ToLower(strings.ReplaceAll(*article.Article.Title, " ", "-"))
}

func (article UpdateArticleDTO) GetSlug() string {
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

	tx, err := repo.DB.BeginTx(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("an error occurred when starting a transaction while attempting to save an article: %w", err)
	}

	var articleId int
	err = tx.QueryRowContext(ctx, query, args...).Scan(&articleId)
	if err != nil {
		switch {
		case err.Error() == "UNIQUE constraint failed: Article.Slug":
			return nil, ErrDuplicateSlug
		default:
			return nil, fmt.Errorf("an error occurred when saving article: %w", err)
		}
	}

	for _, tag := range articleDto.Article.TagList {
		selectTagIdQuery := `SELECT t.TagId FROM Tag t WHERE t.Tag = $1`
		var tagId int

		err = tx.QueryRowContext(ctx, selectTagIdQuery, tag).Scan(&tagId)
		if err == sql.ErrNoRows {

			insertTagQuery := `INSERT INTO Tag (Tag) VALUES ($1) RETURNING TagId`
			err = tx.QueryRowContext(ctx, insertTagQuery, tag).Scan(&tagId)
			if err != nil {
				return nil, fmt.Errorf("an error occurred when attempting to save a tag: %w", err)
			}
		}

		insertArticleTagQuery := `INSERT OR IGNORE INTO ArticleTag (ArticleId, TagId) VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, insertArticleTagQuery, articleId, tagId)
		if err != nil {
			return nil, fmt.Errorf("an error occurred when attempting to tag an article: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("an error occurred when attempting to commit the transaction while saving an article: %w", err)
	}

	article, err := repo.GetArticleBySlug(articleDto.GetSlug(), userId)
	if err != nil {
		return nil, fmt.Errorf("an error occurred when looking up article by slug after saving: %w", err)
	}

	return article, nil
}

func (repo *ArticleRepository) UpdateArticle(articleDto UpdateArticleDTO, articleId, userId int) (*Article, error) {

	query := `UPDATE Article 
			  SET Slug = $1,
			      Title = $2,
				  Description = $3,
				  Body = $4,
				  UpdatedAt = $5
			  WHERE ArticleId = $6 AND UserId = $7`

	now := time.Now().UTC().Format(time.RFC3339Nano)

	args := []any{
		articleDto.GetSlug(),
		*articleDto.Article.Title,
		*articleDto.Article.Description,
		*articleDto.Article.Body,
		now,
		articleId,
		userId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == "UNIQUE constraint failed: Article.Slug":
			return nil, ErrDuplicateSlug
		default:
			return nil, fmt.Errorf("an error occurred when saving article: %w", err)
		}
	}

	article, err := repo.GetArticleBySlug(articleDto.GetSlug(), userId)
	if err != nil {
		return nil, fmt.Errorf("an error occurred when looking up article by slug after saving: %w", err)
	}

	return article, nil
}

func (repo *ArticleRepository) GetArticleBySlug(slug string, userId int) (*Article, error) {
	query := `SELECT
				a.ArticleId,
				a.UserId, 
				a.Title,
				a.Slug,
				a.Description,
				a.Body,
				a.CreatedAt,
				a.UpdatedAt,
				(SELECT EXISTS(SELECT 1 FROM ArticleFavorite af WHERE af.ArticleId=a.ArticleId AND af.UserId=$1)) AS Favorited,
				(SELECT COUNT(*) FROM ArticleFavorite af WHERE af.ArticleId=a.ArticleId) AS FavoritesCount,
				COALESCE((SELECT GROUP_CONCAT(t.Tag, ',') 
				        FROM Tag t 
                        JOIN ArticleTag at ON at.TagId = t.TagId 
                        WHERE at.ArticleId = a.ArticleId), '') AS Tags,
				EXISTS (SELECT 1 FROM Follower WHERE UserId = $1 AND FollowUserId = a.UserId) AS Following,
				u.Username,
				u.Bio,
				u.Image
			  FROM Article a
			  JOIN User u ON a.UserId = u.UserId 
			  WHERE a.Slug = $2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	var article Article
	var author Profile
	var rawTags string
	var createdAt string
	var updatedAt string

	err := repo.DB.QueryRowContext(ctx, query, userId, slug).Scan(
		&article.ArticleId,
		&article.UserId,
		&article.Title,
		&article.Slug,
		&article.Description,
		&article.Body,
		&createdAt,
		&updatedAt,
		&article.Favorited,
		&article.FavoritesCount,
		&rawTags,
		&author.Following,
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

func (repo *ArticleRepository) GetArticles(filters *ArticleFilters, userId int) ([]*BodylessArticle, error) {
	articles := make([]*BodylessArticle, 0)

	query := `SELECT DISTINCT
				a.ArticleId,
				a.UserId, 
				a.Title,
				a.Slug,
				a.Description,
				a.CreatedAt,
				a.UpdatedAt,
				(SELECT EXISTS(SELECT 1 FROM ArticleFavorite af WHERE af.ArticleId=a.ArticleId AND af.UserId=$1)) AS Favorited,
				(SELECT COUNT(*) FROM ArticleFavorite af WHERE af.ArticleId=a.ArticleId) AS FavoritesCount,
				COALESCE((SELECT GROUP_CONCAT(t.Tag, ',') 
						FROM Tag t 
						JOIN ArticleTag at ON at.TagId = t.TagId 
						WHERE at.ArticleId = a.ArticleId), '') AS Tags,
				EXISTS (SELECT 1 FROM Follower WHERE UserId = $1 AND FollowUserId = a.UserId) AS Following,
				u.Username,
				u.Bio,
				u.Image
			FROM Article a
			JOIN User u ON a.UserId = u.UserId
			LEFT JOIN ArticleTag at ON at.ArticleId = a.ArticleId
			LEFT JOIN Tag t ON t.TagId = at.TagId
			LEFT JOIN ArticleFavorite af ON af.ArticleId = a.ArticleId
			LEFT JOIN User favoriter ON favoriter.UserId = af.UserId
			WHERE ($2 IS NULL OR u.Username = $2)
			AND ($3 IS NULL OR favoriter.Username = $3)
			AND ($4 IS NULL OR t.Tag = $4)
			ORDER BY a.ArticleId DESC
			LIMIT $5 OFFSET $6`

	args := []any{userId,
		filters.Author,
		filters.Favorited,
		filters.Tag,
		filters.Limit,
		filters.Offset}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return articles, nil
		default:
			return nil, fmt.Errorf("error querying articles while querying articles: %w", err)
		}
	}

	defer rows.Close()

	for rows.Next() {

		var article BodylessArticle
		var author Profile
		var rawTags string
		var createdAt string
		var updatedAt string

		err = rows.Scan(
			&article.ArticleId,
			&article.UserId,
			&article.Title,
			&article.Slug,
			&article.Description,
			&createdAt,
			&updatedAt,
			&article.Favorited,
			&article.FavoritesCount,
			&rawTags,
			&author.Following,
			&author.Username,
			&author.Bio,
			&author.Image,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning article row: %w", err)
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

		articles = append(articles, &article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows while fetching articles: %w", err)
	}

	return articles, nil
}

func (repo *ArticleRepository) GetFeed(filters *PaginationFilters, userId int) ([]*BodylessArticle, error) {
	articles := make([]*BodylessArticle, 0)

	query := `SELECT
				a.ArticleId,
				a.UserId, 
				a.Title,
				a.Slug,
				a.Description,
				a.CreatedAt,
				a.UpdatedAt,
				(SELECT EXISTS(SELECT 1 FROM ArticleFavorite af WHERE af.ArticleId=a.ArticleId AND af.UserId=$1)) AS Favorited,
				(SELECT COUNT(*) FROM ArticleFavorite af WHERE af.ArticleId=a.ArticleId) AS FavoritesCount,
				COALESCE((SELECT GROUP_CONCAT(t.Tag, ',') 
						FROM Tag t 
						JOIN ArticleTag at ON at.TagId = t.TagId 
						WHERE at.ArticleId = a.ArticleId), '') AS Tags,
				EXISTS (SELECT 1 FROM Follower WHERE UserId = $1 AND FollowUserId = a.UserId) AS Following,
				u.Username,
				u.Bio,
				u.Image
			FROM Article a
			JOIN User u ON a.UserId = u.UserId 
			JOIN Follower f ON a.UserId = f.FollowUserId
			WHERE f.UserId = $1
			ORDER BY a.ArticleId DESC
			LIMIT $2 OFFSET $3`

	args := []any{userId, filters.Limit, filters.Offset}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return articles, nil
		default:
			return nil, fmt.Errorf("error querying articles while constructing feed: %w", err)
		}
	}

	defer rows.Close()

	for rows.Next() {

		var article BodylessArticle
		var author Profile
		var rawTags string
		var createdAt string
		var updatedAt string

		err = rows.Scan(
			&article.ArticleId,
			&article.UserId,
			&article.Title,
			&article.Slug,
			&article.Description,
			&createdAt,
			&updatedAt,
			&article.Favorited,
			&article.FavoritesCount,
			&rawTags,
			&author.Following,
			&author.Username,
			&author.Bio,
			&author.Image,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning article row: %w", err)
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

		articles = append(articles, &article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows while constructing feed: %w", err)
	}

	return articles, nil
}

func (repo *ArticleRepository) DeleteArticle(articleId, userId int) error {
	query := `DELETE FROM Article WHERE ArticleId = $1 AND UserId = $2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, articleId, userId)
	if err != nil {
		return fmt.Errorf("an error occured while trying to delete an article: %w", err)
	}

	return nil
}

func (repo *ArticleRepository) FavoriteArticle(articleId, userId int) error {
	query := `INSERT OR IGNORE INTO ArticleFavorite (ArticleId, UserId) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, articleId, userId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		default:
			return err
		}
	}

	return nil
}

func (repo *ArticleRepository) UnfavoriteArticle(articleId, userId int) error {
	query := `DELETE FROM ArticleFavorite WHERE ArticleId = $1 AND UserId = $2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, articleId, userId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		default:
			return err
		}
	}

	return nil
}
