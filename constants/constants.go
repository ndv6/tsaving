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
	TransferFailed          = "Failed to transfer to main account."
	CannotTransferVaToMain  = "Failed to transfer from virtual account to main account."
	InvalidVA               = "Invalid virtual account number"
	InvalidBalance          = "Insufficient balance"
	CannotParseURLParams    = "Failed to parse URL Params"
	InvalidVaNumber         = "Virtual Account number is invalid"
	FailSqlTransaction      = "Sql Transaction failed to set begin"
	VANotFound              = "Virtual account not found"
	FailToRevertBalance     = "Reverting balance to main account failed"
	TokenExpires            = "Token is already expire, please login to continue"
	TransferToVAFailed      = "Failed transfer to virtual account"
	PasswordRequirement     = "Password Min 6 Character"
	DupeEmailorPhone        = "Unable to Register, Your Phone Number Or Email Has Been Used"
	EmailToken              = "Email Token Failed to Insert"
	AccountFailed           = "Account Number Failed to Insert"
	MailFailed              = "Register Success, but Cannot Send Mail"

	EmailTokenNotFound      = "Can not find requested email"
	VerifyEmailFailed       = "Email fail to be verified with given token"
	UpdateEmailStatusFailed = "Fail to change email status to verified"
	VerifyEmailTokenFailed  = "Unable to verify email token: "
	DeleteEmailTokenFailed  = "Unable to delete verified email"
	SuccessVerifyEmail      = "Email has been successfully verified"
)

//response messages
const (
	AddBalanceVASuccess = "Successfully add balance to your virtual account"
	LoginSucceed        = "Login Succeed"
	RegisterSucceed     = "Register Succeeded"
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
