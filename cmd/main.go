package main

import (
	"net/http"

	conduit "realworld.tayler.io/internal/api"
)

func main() {

	app := conduit.New()

	http.ListenAndServe(":4000", app.Routes())
}
