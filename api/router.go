package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/customers"

	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"

	"github.com/go-chi/chi"
)

func Router(db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)
	chiRouter.Get("/transaction/history", customers.HistoryTransactionHandler(db))

	// Url endpoint not found
	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
