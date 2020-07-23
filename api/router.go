package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/email"

	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"
	"github.com/ndv6/tsaving/api/virtual_accounts"

	"github.com/go-chi/chi"
)

func Router(db *sql.DB) *chi.Mux {
	va := virtual_accounts.NewVAHandler(db)
	chiRouter := chi.NewRouter()
	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)

	// VAC transactions API endpoints
	chiRouter.Post("/vac/to_main", va.VacToMain)
	chiRouter.Get("/vac/list", va.VacList)

	// Url endpoint not found
	// Email verification endpoint
	chiRouter.Post("/email/verify-email-token", email.VerifyEmailToken(db))

	// Not found endpoint
	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
