package constants

// Transaction types
const (
	Deposit = "DEPOSIT_TO_MAIN_ACCOUNT"
)

// Response status
const (
	DepositSuccess = "Deposit completed successfully"
	DepositFailed  = "Failed to make deposit"
)

const (
	Success = "SUCCESS"
	Failed  = "FAILED"
)

// Error messages
const (
	RequestHasInvalidFields = "Mandatory field(s) is missing from request, or is invalid. Please ensure all fields are properly filled."
	CannotReadRequest       = "Cannot read request body"
	CannotParseRequest      = "Unable to parse request"
	UnauthorizedRequest     = "Failed to recognize deposit source. Please use authorized banks to deposit."
	CannotEncodeResponse    = "Failed to encode response."
	InsertFailed            = "Failed to insert to database."
	InitLogFailed           = "Failed to insert log of this transaction."
	TransferToVAFailed      = "Failed transfer to virtual account "
	InvalidBalance          = "Insufficient balance"
	InvalidVA               = "Invalid virtual account number"
)

//response messages
const (
	AddBalanceVASuccess    = "Successfully add balance to your virtual account"
	CannotParseURLParams   = "Failed to parse URL Params"
	TransferFailed         = "Failed to transfer to main account."
	CannotTransferVaToMain = "Failed to transfer from virtual account to main account."
)

// Response messages
const (
	GetListSuccess = "Success to get the list data"
)

// Headers
const (
	ContentType = "Content-Type"
)

// Header types
const (
	Json = "application/json"
)
