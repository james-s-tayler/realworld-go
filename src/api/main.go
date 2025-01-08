package main

import (
	"net/http"
)

type application struct {
}

func main() {
	mux := http.NewServeMux()

	app := &application{}

	mux.HandleFunc("GET /api/user", app.getUserHandler)
	http.ListenAndServe(":4000", mux)
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}
