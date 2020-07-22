package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"

	"github.com/go-chi/chi"
	va "github.com/ndv6/tsaving/api/virtual_accounts"
)

func Router(db *sql.DB) *chi.Mux {
	va := va.NewVAHandler(db)
	chiRouter := chi.NewRouter()
	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)

	// VAC transactions API endpoints
	r.Post("/vac/main", va.vac_to_main(db))

	// Url endpoint not found
	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
