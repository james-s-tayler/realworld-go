package conduit

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"realworld.tayler.io/internal/data"
)

type Application struct {
	logger  *slog.Logger
	domains domains
}

type Config struct {
	DB struct {
		Driver         string
		Dsn            string
		TimeoutSeconds int
	}
}

type domains struct {
	users data.UserRepository
}

type envelope map[string]any

func NewApp(config Config) (*Application, func(), error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger) // so that panics log with slog too

	db, closeDb, err := OpenDB(config, logger)
	if err != nil {
		return nil, nil, err
	}

	app := &Application{
		logger: logger,
		domains: domains{
			users: data.UserRepository{
				DB:             db,
				TimeoutSeconds: config.DB.TimeoutSeconds,
			},
		},
	}

	return app, closeDb, nil
}

func OpenDB(config Config, logger *slog.Logger) (*sql.DB, func(), error) {
	db, err := sql.Open(config.DB.Driver, config.DB.Dsn)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	closeDb := func() {
		if err = db.Close(); err != nil {
			logger.Error(err.Error())
		}
	}

	return db, closeDb, nil
}
