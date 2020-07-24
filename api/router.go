package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/constants"

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

	chiRouter.Use(middleware.Logger)

	// Handler objects initialization
	ph := database.NewPartnerHandler(db)
	ah := database.NewAccountHandler(db)
	ch := customers.NewCustomerHandler(jwt, db)
	va := virtual_accounts.NewVAHandler(jwt, db)

	// Home endpoint
	chiRouter.Get(constants.HomeEndpoint, home.HomeHandler)

	// User onboarding endpoint
	chiRouter.Post(constants.RegisterEndpoint, ch.Create)
	chiRouter.Post(constants.LoginEndpoint, customers.LoginHandler(jwt, db))

	// Virtual Account endpoint
	chiRouter.With(jwt.AuthMiddleware).Post(constants.CreateVirtualAccountEndpoint, vah.Create)
	chiRouter.With(jwt.AuthMiddleware).Put(constants.EditVirtualAccountEndpoint, vah.Edit)

	// VAC transactions API endpoints
	chiRouter.With(jwt.AuthMiddleware).Post(constants.TransferVacToMainAccountEndpoint, va.VacToMain)
	chiRouter.With(jwt.AuthMiddleware).Get(constants.ListAllVacEndpoint, va.VacList)
	chiRouter.With(jwt.AuthMiddleware).Post(constants.DeleteVacEndpoint, va.DeleteVac)
	chiRouter.With(jwt.AuthMiddleware).Put(constants.AddVacBalanceEndpoint, va.AddBalanceVA)

	// Get transaction history
	chiRouter.With(jwt.AuthMiddleware).Get(constants.GetTransactionHistoryEndpoint, ch.HistoryTransactionHandler(db))

	// Email verification endpoint
	chiRouter.Post(constants.VerifyEmailEndpoint, email.VerifyEmailToken(db))

	// Customer Endpoint
	chiRouter.With(jwt.AuthMiddleware).Get(constants.GetUserProfileEndpoint, ch.GetProfile)
	chiRouter.With(jwt.AuthMiddleware).Post(constants.GetUserProfileEndpoint, ch.UpdateProfile)
	chiRouter.With(jwt.AuthMiddleware).Post(constants.GetUserProfileEndpoint, ch.UpdatePhoto)

	// Main account transactions endpoint
	chiRouter.Post(constants.DepositEndpoint, customers.DepositToMainAccount(ph, ah))

	// Url endpoint not found
	chiRouter.NotFound(not_found.NotFoundHandler)

	return chiRouter
}
