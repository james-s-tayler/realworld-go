package conduit

import (
	"log/slog"
	"os"
)

type Application struct {
	logger *slog.Logger
}

func New() *Application {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger) // so that panics log with slog too

	return &Application{
		logger: logger,
	}
}
