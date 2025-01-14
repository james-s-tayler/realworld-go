package main

import (
	"fmt"
	"net/http"

	conduit "realworld.tayler.io/internal/api"
)

func main() {
	config := conduit.Config{
		DB: struct {
			Driver         string
			Dsn            string
			TimeoutSeconds int
		}{
			Driver:         "sqlite3",
			Dsn:            "file:conduit.db?mode=rwc&cache=shared",
			TimeoutSeconds: 30,
		},
	}
	app, closeDb, err := conduit.NewApp(config)
	if err != nil {
		fmt.Printf("Error starting the application: %v", err.Error())
	}

	defer closeDb()

	http.ListenAndServe(":4000", app.Routes())
}
