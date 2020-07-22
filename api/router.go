package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/email"

	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"

	"github.com/go-chi/chi"
)

func Router(db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)

	// Email verification endpoint
	chiRouter.Post("/email/verify-email-token", email.VerifyEmailToken(db))

	// Not found endpoint
	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
