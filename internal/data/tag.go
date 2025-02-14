package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type TagRepository struct {
	DB             *sql.DB
	TimeoutSeconds int
}

func (repo *TagRepository) GetAllTags() ([]string, error) {
	tags := make([]string, 0)

	query := `SELECT DISTINCT tag FROM Tag ORDER BY TagId ASC`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repo.TimeoutSeconds)*time.Second)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return tags, nil
		default:
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var tag string

		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
