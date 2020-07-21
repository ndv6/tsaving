package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
)

func Router(db *sql.DB) http.Handler {
	va := virtual_accounts.NewVAHandler(db)
	r := chi.NewRouter()
	// r.Post("/login", LoginHandler(jwt, db))
	// r.Get("/", HomeHandler)
	// r.With().Get("/customers", ch.List)          //ini get jenisnya.
	// r.With().Post("/customer/create", ch.Create) // kalau ini tesnya POST.
	// r.With().Get("/customer/{id}", ch.Get)       // kalau di chi bisa langsung di define expect apa.
	r.Post("/vac/main", va.vac_to_main(db))
	// r.NotFound(NotFound)

	return r
}
