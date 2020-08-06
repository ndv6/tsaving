package api

import (
	"database/sql"

	"github.com/go-chi/chi/middleware"

	"github.com/ndv6/tsaving/database"

	"github.com/ndv6/tsaving/api/admin"
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
	chiRouter.Use(middleware.Recoverer)

	// Handler objects initialization
	ph := database.NewPartnerHandler(db)
	ah := database.NewAccountHandler(db)
	ch := customers.NewCustomerHandler(jwt, db)
	va := virtual_accounts.NewVAHandler(jwt, db) // David, Jocelyn, Joseph , Azizah
	eh := database.NewEmailHandler(db)           // Joseph
	la := admin.NewLogAdminHandler(jwt, db)
	// da := customers.NewDashboardHandler(jwt)

	// Home endpoint
	chiRouter.Get("/", home.HomeHandler)

	// Registration Endpoint
	chiRouter.Post("/register", ch.Create)                        //Caesar
	chiRouter.Post("/verify-account", email.VerifyEmailToken(eh)) //Joseph
	chiRouter.Post("/get-token", email.GetEmailToken(eh))         // Yuly

	// Login Endpoint
	chiRouter.Post("/login", customers.LoginHandler(jwt, db)) //Caesar

	// Customer Endpoint
	chiRouter.With(jwt.AuthMiddleware).Get("/me/profile", ch.GetProfile)                                 //Andreas
	chiRouter.With(jwt.AuthMiddleware).Put("/me/update", ch.UpdateProfile)                               //Andreas
	chiRouter.With(jwt.AuthMiddleware).Patch("/me/update-photo", ch.UpdatePhoto)                         //Andreas
	chiRouter.With(jwt.AuthMiddleware).Patch("/me/update-password", ch.UpdatePassword)                   //Andreas
	chiRouter.With(jwt.ValidateAccount).Post("/me/deposit", customers.DepositToMainAccount(ph, ah))      //Vici
	chiRouter.With(jwt.AuthMiddleware).With(jwt.ValidateAccount).Put("/me/transfer-va", va.AddBalanceVA) //David
	chiRouter.With(jwt.AuthMiddleware).Get("/me/dashboard", ch.GetDashboardData(db))                     //David

	// Virtual Account Endpoint
	chiRouter.With(jwt.AuthMiddleware).Get("/me/va", va.VacList)                                                     //Jocelyn
	chiRouter.With(jwt.AuthMiddleware).With(jwt.ValidateAccount).Post("/me/va/create", va.Create)                    //Azizah
	chiRouter.With(jwt.AuthMiddleware).With(jwt.ValidateAccount).Put("/me/va/{va_num}/update", va.Update)            //Azizah
	chiRouter.With(jwt.AuthMiddleware).With(jwt.ValidateAccount).Post("/me/va/{va_num}/transfer-main", va.VacToMain) //Jocelyn
	chiRouter.With(jwt.AuthMiddleware).With(jwt.ValidateAccount).Delete("/me/va/{va_num}", va.DeleteVac)             //Joseph

	// History Endpoint -- Yuly Haruka
	chiRouter.With(jwt.AuthMiddleware).Get("/me/transaction/{page}", ch.HistoryTransactionHandler(db)) //Yuly

	// Log Admin
	chiRouter.With(jwt.AuthMiddleware).Get("/admin/log/{page}", la.Get)
	chiRouter.With(jwt.AuthMiddleware).Post("/admin/log/insert", la.Insert)

	// Not Found Endpoint
	chiRouter.NotFound(not_found.NotFoundHandler) // Joseph

	return chiRouter
}
