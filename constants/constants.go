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

// Error messages
const (
	RequestHasInvalidFields = "Mandatory field(s) is missing from request, or is invalid. Please ensure all fields are properly filled."
	CannotReadRequest       = "Cannot read request body"
	CannotParseRequest      = "Unable to parse request"
	UnauthorizedRequest     = "Failed to recognize deposit source. Please use authorized banks to deposit."
	CannotEncodeResponse    = "Failed to encode response."
	InsertFailed            = "Failed to insert to database."
	InitLogFailed           = "Failed to insert log of this transaction."
)
