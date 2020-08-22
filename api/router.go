package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/ndv6/tsaving/api/static"
	"github.com/ndv6/tsaving/constants"

	"github.com/ndv6/tsaving/models"

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
	th := models.NewTokenHandler(db)
	admDb := database.NewAdminDatabaseHandler(db)
	ph := database.NewPartnerHandler(db)
	ah := database.NewAccountHandler(db)
	ch := customers.NewCustomerHandler(jwt, db)
	va := virtual_accounts.NewVAHandler(jwt, db) // David, Jocelyn, Joseph , Azizah
	eh := database.NewEmailHandler(db)           // Joseph
	adm := admin.NewAdminHandler(jwt, db)        // Azizah
	la := admin.NewLogAdminHandler(jwt, db)
	pph := customers.NewPictureHandler(jwt)

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
	chiRouter.With(jwt.AuthMiddleware).Patch("/me/update-photo", pph.V2UpdatePhoto)                      //Andreas
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
		r.With(jwt.AuthAdminMiddleware).Put("/customers/edit", adm.EditCustomerData(admDb, th))    //Vici

		// transaction log
		r.Route("/transactions", func(r chi.Router) {
			r.With(jwt.AuthAdminMiddleware).Get("/", adm.TransactionHistoryHandler)                                              // Azizah
			r.With(jwt.AuthAdminMiddleware).Get("/{accNum}/{page}", adm.TransactionHistoryHandler)                               // Yuly
			r.With(jwt.AuthAdminMiddleware).Get("/{accNum}/{day}-{month}-{year}/{page}", adm.TransactionHistoryHandler)          // Yuly
			r.With(jwt.AuthAdminMiddleware).Get("/{accNum}/{search}/{page}", adm.TransactionHistoryHandler)                      // Yuly
			r.With(jwt.AuthAdminMiddleware).Get("/{accNum}/{day}-{month}-{year}/{search}/{page}", adm.TransactionHistoryHandler) // Yuly
			r.With(jwt.AuthAdminMiddleware).Get("/list/{page}", adm.TransactionHistoryAll)                                       // Azizah
			r.With(jwt.AuthAdminMiddleware).Get("/list/d/{date}/{page}", adm.TransactionHistoryAll)                              // Azizah
			r.With(jwt.AuthAdminMiddleware).Get("/list/a/{search}/{page}", adm.TransactionHistoryAll)                            // Azizah
			r.With(jwt.AuthAdminMiddleware).Get("/list/{search}/{date}/{page}", adm.TransactionHistoryAll)                       // Azizah
		})

		// Log Admin
		r.Route("/log", func(r chi.Router) {
			r.With(jwt.AuthAdminMiddleware).Get("/{page}", la.Get)                            //Jocelyn
			r.With(jwt.AuthAdminMiddleware).Post("/insert", la.Insert)                        //Jocelyn
			r.With(jwt.AuthAdminMiddleware).Get("/d/{date}/{page}", la.GetFilteredLog)        //Azizah
			r.With(jwt.AuthAdminMiddleware).Get("/u/{search}/{page}", la.GetFilteredLog)      //Azizah
			r.With(jwt.AuthAdminMiddleware).Get("/{search}/{date}/{page}", la.GetFilteredLog) //Azizah
		})

		// admin dashboard
		r.With(jwt.AuthAdminMiddleware).Get("/dashboard", adm.GetDashboard()) // Joseph

		// va list for admin
		r.With(jwt.AuthAdminMiddleware).Get("/va/{cust_id}/{page}", va.VacListAdmin)
		r.With(jwt.AuthAdminMiddleware).Get("/va/{cust_id}/{color}/{page}", va.VacListAdminFilter)

		// get token for resend email
		r.With(jwt.AuthAdminMiddleware).Post("/get-token", email.GetEmailToken(eh))
	})

	FileServer(chiRouter, "/static", http.Dir("./static"), jwt)

	// Not Found Endpoint
	chiRouter.NotFound(not_found.NotFoundHandler) // Joseph

	return chiRouter
}

// since go chi doesnt provide file server, they make a workaround for that
// https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func FileServer(r chi.Router, path string, root http.FileSystem, jwt *tokens.JWT) {
	if strings.ContainsAny(path, "{}*") {
		panic(constants.FileServerConfigError)
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	fmt.Println("Server is serving static files at: host/static")
	r.With(jwt.AuthMiddleware).Get(path, static.StaticHandler(root))
}
