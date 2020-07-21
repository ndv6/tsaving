package virtual_accounts

import (
	"database/sql"
)

type VAHandler struct {
	db *sql.DB
}

func NewVAHanlder(db *sql.DB) *VAHandler {
	return &VAHandler{db}
}

func vac_to_main() {

	//ambil input dari jsonnya (no rek VAC dan saldo input)

	//cek kalau no rek input bener ga
	helper.check_rekening()

	//cek input apakah melebihi saldo?

	//ambil id customer dari token

	// panggil fungsi update (saldo vac - input) (saldo utama + input) by id cust.
	return

}

//function di model
