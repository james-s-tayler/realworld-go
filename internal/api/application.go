package conduit

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"realworld.tayler.io/internal/data"
)

type Application struct {
	logger       *slog.Logger
	domains      domains
	tokenService data.ITokenService
}

type Config struct {
	DB struct {
		Driver         string
		Dsn            string
		TimeoutSeconds int
	}
	JWT struct {
		SecretKey []byte
	}
}

type domains struct {
	users    data.UserRepository
	articles data.ArticleRepository
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
			articles: data.ArticleRepository{
				DB:             db,
				TimeoutSeconds: config.DB.TimeoutSeconds,
			},
		},
		tokenService: data.JwtTokenService{
			SecretKey: config.JWT.SecretKey,
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

	// Enable foreign keys - this is required on every connection for sqlite
	// since sqlite doesn't enforce foreign keys by default for backwards compatibility reasons
	_, err = db.ExecContext(ctx, "PRAGMA foreign_keys = ON;")
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

func (app *Application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Use the json.MarshalIndent() function so that whitespace is added to the encoded
	// JSON. Here we use no line prefix ("") and tab indents ("\t") for each element.
	// There's a small perf hit for this, but it's neglible.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// this makes it easier to read in terminal applications
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	/*
	* the "json" package in Go has some flaws and a v2 is in discussion here:
	* https://github.com/golang/go/discussions/63397
	 */

	// this limits the max size of the request body to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		// if there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unsharmalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body ontains badly=formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unsharmalTypeError):
			if unsharmalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unsharmalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unsharmalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message. Note that there's an open
		// issue at https://github.com/golang/go/issues/29035 regarding turning this
		// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
