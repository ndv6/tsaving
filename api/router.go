package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/tokens"

	"github.com/ndv6/tsaving/api/customers"
	"github.com/ndv6/tsaving/api/email"
	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"

	"github.com/go-chi/chi"
)

func Router(jwt *tokens.JWT, db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	ch := customers.NewCustomerHandler(db, jwt)
	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)

	// Get transaction history
	chiRouter.Get("/transaction/history", ch.HistoryTransactionHandler(db))

	// Email verification endpoint
	chiRouter.Post("/email/verify-email-token", email.VerifyEmailToken(db))

	// Not found endpoint
	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
