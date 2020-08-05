package admin

type LogAdminHandler struct {
	jwt *tokens.JWT
	db  *sql.DB
}

func NewLogAdminHandler(jwt *tokens.JWT, db *sql.DB) *VAHandler {
	return &LogAdminHandler{jwt, db}
}

func (la *LogAdminHandler) VacToMain(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set(constants.ContentType, constants.Json)
	token := va.jwt.GetToken(r)
	//ambil input dari jsonnya (no rek VAC dan saldo input)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.CannotReadRequest)
		return
	}

	var LogAdmin models.LogAdmin
	err = json.Unmarshal(b, &LogAdmin)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.CannotParseRequest)
		return
	}

	//cek admin ?
	UsernameAdmin := chi.URLParam(r,"username") //ini bisa ambil dari token
	err = database.CheckAdmin(la.db,UsernameAdmin)
	if err != nil {
		helper.HTTPError(w, http.StatusBadRequest, constants.InvalidAdmin)
		return
	}

	database.GetLogAdmin(la.db, UsernameAdmin)


	//get no rekening by rekening vac
	AccountNumber, _ := database.GetAccountByVA(va.db, vaNum)

	//update balance at both accounts
	err = database.UpdateVacToMain(va.db, VirAcc.BalanceChange, vaNum, AccountNumber)
	if err != nil {
		helper.HTTPError(w, http.StatusOK, err.Error())
		return
	}

	_, res, err := helpers.NewResponseBuilder(w, true, fmt.Sprintf("successfully move balance to your main account : %v", VirAcc.BalanceChange), nil)
	if err != nil {
		helpers.HTTPError(w, http.StatusBadRequest, constants.CannotEncodeResponse)
		return
	}
	fmt.Fprint(w, string(res))

	return

}
