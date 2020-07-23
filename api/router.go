package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/customers"
	"github.com/ndv6/tsaving/api/email"
	"github.com/ndv6/tsaving/api/not_found"
	"github.com/ndv6/tsaving/tokens"

	"github.com/ndv6/tsaving/api/home"

	"github.com/go-chi/chi"
	"github.com/ndv6/tsaving/api/virtual_accounts"
)

func Router(jwt *tokens.JWT, db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	chiRouter.Get("/", home.HomeHandler)
	va := virtual_accounts.NewVAHandler(jwt, db)
	chiRouter.With(jwt.AuthMiddleware).Put("/vac/add_balance_vac", va.AddBalanceVA)

	ch := customers.NewCustomerHandler(jwt, db)
	chiRouter.Get("/", home.HomeHandler)
	chiRouter.Post("/register", ch.Create)
	chiRouter.Post("/login", customers.LoginHandler(jwt, db))

	chiRouter.Post("/email/verify-email-token", email.VerifyEmailToken(db))

	chiRouter.NotFound(not_found.NotFoundHandler)
	return chiRouter
}
