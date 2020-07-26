package api

import (
	"database/sql"

	"github.com/go-chi/chi/middleware"

	"github.com/ndv6/tsaving/database"

	"github.com/ndv6/tsaving/api/customers"
	"github.com/ndv6/tsaving/api/email"
	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"
	"github.com/ndv6/tsaving/api/virtual_accounts"
	"github.com/ndv6/tsaving/tokens"

	"github.com/go-chi/chi"
)

func Router(jwt *tokens.JWT, db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.Logger)

	// Handler objects initialization
	ph := database.NewPartnerHandler(db)
	ah := database.NewAccountHandler(db)
	ch := customers.NewCustomerHandler(jwt, db)
	va := virtual_accounts.NewVAHandler(jwt, db)

	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)
	chiRouter.Post("/register", ch.Create)
	chiRouter.Post("/login", customers.LoginHandler(jwt, db))

	// VAC transactions API endpoints
	chiRouter.With(jwt.AuthMiddleware).Post("/vac/to_main", va.VacToMain)
	chiRouter.With(jwt.AuthMiddleware).Get("/vac/list", va.VacList)
	chiRouter.With(jwt.AuthMiddleware).Post("/vac/delete-vac", va.DeleteVac)

	// Registration Endpoint
	chiRouter.Post("/register", ch.Create)
	chiRouter.Post("/verify-account", email.VerifyEmailToken(db))

	// Login Endpoint
	chiRouter.Post("/login", customers.LoginHandler(jwt, db))

	// Customer Endpoint
	chiRouter.With(jwt.AuthMiddleware).Get("/me/profile", ch.GetProfile)
	chiRouter.With(jwt.AuthMiddleware).Put("/me/update", ch.UpdateProfile)
	chiRouter.With(jwt.AuthMiddleware).Patch("/me/update-photo", ch.UpdatePhoto)
	chiRouter.Post("/me/deposit", customers.DepositToMainAccount(ph, ah))
	chiRouter.With(jwt.AuthMiddleware).Put("/me/transfer-va", va.AddBalanceVA)

	// Virtual Account Endpoint
	chiRouter.With(jwt.AuthMiddleware).Get("/me/va", va.VacList)
	chiRouter.With(jwt.AuthMiddleware).Post("/me/va/create", va.Create)
	chiRouter.With(jwt.AuthMiddleware).Put("/me/va/{va_num}/update", va.Update)
	chiRouter.With(jwt.AuthMiddleware).Post("/me/va/{va_num}/transfer-main", va.VacToMain)
	chiRouter.With(jwt.AuthMiddleware).Delete("/me/va/{va_num}", va.DeleteVac)

	// History Endpoint
	chiRouter.With(jwt.AuthMiddleware).Get("/me/transaction/{page}", ch.HistoryTransactionHandler(db))

	// Not Found Endpoint
	chiRouter.NotFound(not_found.NotFoundHandler)

	return chiRouter
}
