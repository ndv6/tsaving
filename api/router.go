package api

import (
	"database/sql"

	"github.com/go-chi/chi/middleware"

	"github.com/ndv6/tsaving/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/ndv6/tsaving/api/admin"
	"github.com/ndv6/tsaving/api/customers"
	"github.com/ndv6/tsaving/api/email"
	"github.com/ndv6/tsaving/api/home"
	"github.com/ndv6/tsaving/api/not_found"
	"github.com/ndv6/tsaving/api/virtual_accounts"
	"github.com/ndv6/tsaving/tokens"
)

func Router(jwt *tokens.JWT, db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.Logger)
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	// Handler objects initialization
	ph := database.NewPartnerHandler(db)
	ah := database.NewAccountHandler(db)
	ch := customers.NewCustomerHandler(jwt, db)
	va := virtual_accounts.NewVAHandler(jwt, db) // David, Jocelyn, Joseph , Azizah
	eh := database.NewEmailHandler(db)           // Joseph
	adm := admin.NewAdminHandler(db)             // Azizah
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

	// History Endpoint
	chiRouter.With(jwt.AuthMiddleware).Get("/me/transaction/{page}", ch.HistoryTransactionHandler(db)) //Yuly

	chiRouter.Route("/v2", func(r chi.Router) {
		// login
		r.Post("/login", admin.LoginAdminHandler(jwt, db)) //Caesar

		// customer details
		r.With(jwt.AuthAdminMiddleware).Post("/customers/list/{page}", ch.GetListCustomers)        //David
		r.With(jwt.AuthAdminMiddleware).Get("/customers/cards/{account_num}", ch.GetCardCustomers) //Caesar
		r.With(jwt.AuthAdminMiddleware).Get("/customers/{cust_id}", ch.GetProfileforAdmin)         //Caesar
		r.With(jwt.AuthAdminMiddleware).Post("/customers/delete", ch.SoftDelete)                   //Jocelyn

		// transaction log
		r.Route("/transactions", func(r chi.Router) {
			r.With(jwt.AuthAdminMiddleware).Get("/", adm.TransactionHistoryHandler)                  // Azizah
			r.With(jwt.AuthAdminMiddleware).Get("/{accNum}", adm.TransactionHistoryHandler)          // Yuly
			r.With(jwt.AuthAdminMiddleware).Get("/{accNum}/{search}", adm.TransactionHistoryHandler) // Yuly
		})

		// Log Admin
		r.With(jwt.AuthAdminMiddleware).Get("/log/{page}", la.Get)     //Jocelyn
		r.With(jwt.AuthAdminMiddleware).Post("/log/insert", la.Insert) //Jocelyn

		// admin dashboard
		r.With(jwt.AuthAdminMiddleware).Get("/dashboard", adm.GetDashboard())

		// va list for admin
		r.With(jwt.AuthAdminMiddleware).Get("/va/{cust_id}", va.VacListAdmin)
	})

	// Not Found Endpoint
	chiRouter.NotFound(not_found.NotFoundHandler) // Joseph

	return chiRouter
}
