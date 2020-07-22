package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ndv6/tsaving/api/virtual_accounts"
)

func Router(db *sql.DB) *chi.Mux {
	chiRouter := chi.NewRouter()
	chiRouter.Get("/", HomeHandler)
	chiRouter.NotFound(NotFoundHandler)
	va := virtual_accounts.NewVAHandler(db)
	chiRouter.Post("/vac/add_balance_vac", va.AddBalanceVA)
	chiRouter.Post("/", HomeHandler)
	return chiRouter
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Tsaving")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Seems like %v is not available or does not exist", r.URL)
}
