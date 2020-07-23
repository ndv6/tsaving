package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/customers"
	"github.com/ndv6/tsaving/api/email"
	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"
	"github.com/ndv6/tsaving/tokens"

	"github.com/go-chi/chi"
)

func Router(jwt *tokens.JWT, db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	ch := customers.NewCustomerHandler(jwt, db)
	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)
	chiRouter.Post("/register", ch.Create)
	chiRouter.Post("/login", customers.LoginHandler(jwt, db))

	// Get transaction history
	chiRouter.With(jwt.AuthMiddleware).Get("/transaction/history", ch.HistoryTransactionHandler(db))

	// Email verification endpoint
	chiRouter.Post("/email/verify-email-token", email.VerifyEmailToken(db))

	// Customer Endpoint
	chiRouter.With(jwt.AuthMiddleware).Get("/customers/getprofile", ch.GetProfile)
	chiRouter.With(jwt.AuthMiddleware).Post("/customers/updateprofile", ch.UpdateProfile)
	chiRouter.With(jwt.AuthMiddleware).Post("/customers/updatephoto", ch.UpdatePhoto)

	// Not found endpoint
	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
