package api

import (
	"database/sql"

	"github.com/go-chi/chi/middleware"

	"github.com/ndv6/tsaving/api/customers"
	"github.com/ndv6/tsaving/api/email"

	"github.com/ndv6/tsaving/database"

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

	// Virtual Account endpoint
	chiRouter.With(jwt.AuthMiddleware).Post("/me/va/create", va.Create)
	chiRouter.With(jwt.AuthMiddleware).Put("/me/va/{va_num}/update", va.Update)

	// VAC transactions API endpoints
	chiRouter.With(jwt.AuthMiddleware).Post("/vac/to_main", va.VacToMain)
	chiRouter.With(jwt.AuthMiddleware).Get("/vac/list", va.VacList)
	chiRouter.With(jwt.AuthMiddleware).Post("/vac/delete-vac", va.DeleteVac)

	// Url endpoint not found
	// Get transaction history
	chiRouter.With(jwt.AuthMiddleware).Get("/transaction/history", ch.HistoryTransactionHandler(db))

	// Email verification endpoint
	chiRouter.Post("/email/verify-email-token", email.VerifyEmailToken(db))

	// Customer Endpoint
	chiRouter.With(jwt.AuthMiddleware).Get("/customers/getprofile", ch.GetProfile)
	chiRouter.With(jwt.AuthMiddleware).Post("/customers/updateprofile", ch.UpdateProfile)
	chiRouter.With(jwt.AuthMiddleware).Post("/customers/updatephoto", ch.UpdatePhoto)

	// Main account transactions endpoint
	chiRouter.Post("/deposit", customers.DepositToMainAccount(ph, ah))

	// Url endpoint not found
	chiRouter.NotFound(not_found.NotFoundHandler)

	return chiRouter
}
