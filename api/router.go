package api

import (
	"database/sql"

	"github.com/polipopoliko/ndv6/tsaving/api/home/homepage"
	"github.com/polipopoliko/ndv6/tsaving/api/not_found/not_found_page"

	"github.com/go-chi/chi"
)

func Router(db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	chiRouter.Get("/", homepage.HomeHandler)
	chiRouter.NotFound(not_found_page.NotFoundHandler)
	return chiRouter
}
