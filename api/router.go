package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"
	"github.com/ndv6/tsaving/api/virtual_accounts"

	"github.com/go-chi/chi"
)

func Router(db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	vah := virtual_accounts.NewVirtualAccHandler(db)
	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)

	// Virtual Account endpoint
	chiRouter.Post("/virtualaccount/create", vah.Create)
	chiRouter.Put("/virtualaccount/edit/{vaNumber}", vah.Edit)

	// Url endpoint not found
	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
