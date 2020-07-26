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
	vah := virtual_accounts.NewVAHandler(jwt, db)

	// to log incoming requests
	chiRouter.Use(middleware.Logger)

	// Handler objects initialization
	ph := database.NewPartnerHandler(db)
	ah := database.NewAccountHandler(db)
	ch := customers.NewCustomerHandler(jwt, db)
	va := virtual_accounts.NewVAHandler(jwt, db)
	eh := database.EmailHandler{Db: db}

	// Home Endpoint
	chiRouter.Get("/", home.HomeHandler)

	// Virtual Account endpoint
	chiRouter.With(jwt.AuthMiddleware).Post("/virtualaccount/create", vah.Create)
	chiRouter.With(jwt.AuthMiddleware).Put("/virtualaccount/edit", vah.Edit)
	// VAC transactions API endpoints
	chiRouter.With(jwt.AuthMiddleware).Post("/me/va/{va_num}/transfer-to-main", va.VacToMain)
	chiRouter.With(jwt.AuthMiddleware).Get("/me/va/list", va.VacList)
	chiRouter.With(jwt.AuthMiddleware).Post("/vac/delete-vac", va.DeleteVac)

	// Get transaction history
	chiRouter.With(jwt.AuthMiddleware).Get("/transaction/history/{page}", ch.HistoryTransactionHandler(db))
	// Registration Endpoint
	chiRouter.Post("/register", ch.Create)
	chiRouter.Post("/verify-account", email.VerifyEmailToken(eh))

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
	chiRouter.With(jwt.AuthMiddleware).Post("/me/va/create", vah.Create)
	chiRouter.With(jwt.AuthMiddleware).Put("/me/va/{va_num}", vah.Edit)
	chiRouter.With(jwt.AuthMiddleware).Post("/me/va/{va_num}/transfer-main", va.VacToMain)
	chiRouter.With(jwt.AuthMiddleware).Delete("/me/va/{va_num}", va.DeleteVac)

	// History Endpoint -- Yuly Haruka
	chiRouter.With(jwt.AuthMiddleware).Get("/me/transaction/{page}", ch.HistoryTransactionHandler(db))

	// Not Found Endpoint
	chiRouter.NotFound(not_found.NotFoundHandler)

	return chiRouter
}
