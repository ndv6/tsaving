package virtual_accounts

import "net/http"

type VirtualAccount struct {
	Rekening_VAC string `json:"rekening_vac"`
	Saldo_VAC    string `json:"saldo_vac"`
}

func vac_to_main(db) {
	return func(w http.ResponseWriter, r *http.Request) {
		var va VirtualAccount

		//ambil input dari jsonnya (no rek VAC dan saldo input)

		//cek kalau no rek input bener ga

		//cek input apakah melebihi saldo?

		//ambil id customer dari token

		// panggil fungsi update (saldo vac - input) (saldo utama + input) by id cust.
		return
	}
}

//function di model
