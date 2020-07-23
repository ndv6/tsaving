package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/not_found"

	"github.com/ndv6/tsaving/api/home"

	"github.com/go-chi/chi"
	"github.com/ndv6/tsaving/api/virtual_accounts"
)

func Router(db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	chiRouter.Get("/", home.HomeHandler)
	va := virtual_accounts.NewVAHandler(db)
	chiRouter.Put("/vac/add_balance_vac", va.AddBalanceVA)
	// ch := customers.New
	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
