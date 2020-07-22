package api

import (
	"database/sql"

	"github.com/ndv6/tsaving/api/email"

	"github.com/ndv6/tsaving/database"

	"github.com/ndv6/tsaving/api/customers"
	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"
	"github.com/ndv6/tsaving/tokens"

	"github.com/go-chi/chi"
)

func Router(db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()

	// Handler objects initialization
	ph := database.NewPartnerHandler(db)
	ah := database.NewAccountHandler(db)
<<<<<<< HEAD
=======
	ch := customers.NewCustomerHandler(jwt, db)
>>>>>>> 4505c76... feat: api to deposit to main account

	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)

	// Email verification endpoint
	chiRouter.Post("/email/verify-email-token", email.VerifyEmailToken(db))

	// Main account transactions endpoint
	chiRouter.Post("/deposit", customers.DepositToMainAccount(ph, ah))

	// Url endpoint not found
	chiRouter.NotFound(not_found.NotFoundHandler)

	return chiRouter
}
