package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/email"

	"github.com/ndv6/tsaving/api/customers"
	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"
	"github.com/ndv6/tsaving/api/virtual_accounts"
	"github.com/ndv6/tsaving/tokens"

	"github.com/go-chi/chi"
)

func Router(jwt *tokens.JWT, db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	ch := customers.NewCustomerHandler(jwt, db)
	va := virtual_accounts.NewVAHandler(jwt, db)
	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)
	chiRouter.Post("/register", ch.Create)
	chiRouter.Post("/login", customers.LoginHandler(jwt, db))

	// VAC transactions API endpoints
	chiRouter.With(jwt.AuthMiddleware).Post("/vac/to_main", va.VacToMain)
	chiRouter.With(jwt.AuthMiddleware).Get("/vac/list", va.VacList)

	// Url endpoint not found
	// Email verification endpoint
	chiRouter.Post("/email/verify-email-token", email.VerifyEmailToken(db))

	// Not found endpoint
	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
