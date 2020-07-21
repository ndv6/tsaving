package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func Router(db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	chiRouter.Get("/", HomeHandler)
	chiRouter.NotFound(NotFoundHandler)
	return chiRouter
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Tsaving")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Seems like %v is not available or does not exist", r.URL)
}
