package constants

// Transaction types
const (
	Deposit                  = "DEPOSIT_TO_MAIN_ACCOUNT"
	TransferToVirtualAccount = "MAIN_TO_VA"
	TransferToMainAccount    = "VA_TO_MAIN"
)

// Admin Activity Types (for Admin Activity Log)
const (
	EditCustomerData = "EDIT"
	DELETE           = "DELETE"
)

const (
	Success = "SUCCESS"
	Failed  = "FAILED"
)

// Error messages
const (
	GenerateEmailTokenFailed = "Failed to generate new e-mail token. Please resend verification e-mail to customer."
	EditSuccessMailNotSent   = "Edit customer data success, but verification e-mail was not sent. Please resend e-mail ro customer."
	RequestHasInvalidFields  = "Mandatory field(s) is missing from request, or is invalid. Please ensure all fields are properly filled."
	CannotReadRequest        = "Cannot read request body"
	CannotParseRequest       = "Unable to parse request"
	UnauthorizedRequest      = "Failed to recognize deposit source. Please use authorized banks to deposit."
	CannotEncodeResponse     = "Failed to encode response."
	InsertFailed             = "Failed to insert to database."
	UpdateFailed             = "Failed to update to database."
	InitLogFailed            = "Failed to insert log of this transaction."
	CannotParseURLParams     = "Failed to parse URL Params"
	InvalidUrlParams         = "Url params not valid"
	FailToLogin              = "Unable to login using given credential"
	MinimumPassword          = "The password must have at least 6 characters"
	EmailTaken               = "The email you entered is already taken"
	InvalidEmail             = "Invalid email address format"
	PhoneTaken               = "The phone number you entered is already taken"
	TransferToVAFailed       = "Failed transfer to virtual account "
	InvalidBalance           = "Insufficient balance"
	TransferFailed           = "Failed to transfer to main account."
	CannotTransferVaToMain   = "Failed to transfer from virtual account to main account."
	InvalidVA                = "Invalid virtual account number"
	InvalidVaNumber          = "Virtual Account number is invalid"
	FailSqlTransaction       = "Sql Transaction failed to set begin"
	VANotFound               = "Virtual account not found"
	FailToRevertBalance      = "Reverting balance to main account failed"
	TokenExpires             = "Token is already expire, please login to continue"
	PasswordRequirement      = "Password Min 6 Character"
	DupeEmailorPhone         = "Unable to Register, Your Phone Number Or Email Has Been Used"
	EmailToken               = "Email Token Failed to Insert"
	AccountFailed            = "Account Number Failed to Insert"
	MailFailed               = "Register Success, but Cannot Send Mail"
	EmailTokenNotFound       = "Can not find requested email"
	VerifyEmailFailed        = "Email fail to be verified with given token"
	UpdateEmailStatusFailed  = "Fail to change email status to verified"
	VerifyEmailTokenFailed   = "Unable to verify email token: "
	DeleteEmailTokenFailed   = "Unable to delete verified email"
	SuccessVerifyEmail       = "Email has been successfully verified"
	LogAdminFailed           = "Failed to get admin log."
	InsertAdminLogFailed     = "Insert Admin Log Failed"
	CustomersNotFound        = "List Customers Not Found"
	GetAllTransactionFailed  = "Failed to get all transaction"
	InvalidAccountNumber     = "Invalid account number"
	SoftDeleteCustFailed     = "Soft delete customer failed"
	LogAdminCannotBeNull     = "Username or Filter cannot be null"
	DateCannotBeNull         = "Date cannot be null"
	UsernameCannotBeNull     = "Username cannot be null"
	Base64DecodeFailed       = "Fail to decode input as base64"
)

// Response messages
const (
	DepositSuccess           = "Deposit completed successfully"
	DepositFailed            = "Failed to make deposit"
	EditCustomerDataSuccess  = "Customer data updated successfully."
	AddBalanceVASuccess      = "Successfully add balance to your virtual account"
	LoginSucceed             = "Login Succeed"
	RegisterSucceed          = "Register Succeeded"
	GetListSuccess           = "Success to get the list data"
	GetAllTransactionSuccess = "Success Get All Transactions"
	GetProfilSuccess         = "Get profile success"
	GetCardSuccess           = "Get Cards Success"
	UpdateProfileSuccess     = "Update profile success"
	UpdatePasswordSuccess    = "Update password success"
	UpdatePhotoSuccess       = "Update photo success"
	SuccessGetToken          = "Success to get the token"
	AddLogAdminSuccess       = "Add admin log success"
	GetLogAdminSuccess       = "Get admin log success"
	SuccessSoftDelete        = "Customer deleted successfully."
)

// Static Handler
const (
	FileServerConfigError = "FileServer does not permit any URL parameters."
	UpdateImageFailed     = "Fail to update image"
	UnsupportedImageType  = "Uploaded image extension should be either JPEG or PNG"
	StaticPath            = "/static"
	StaticImagePath       = StaticPath + "/images"
)

// Headers
const (
	ContentType = "Content-Type"
)

// Header types
const (
	Json = "application/json"
)

// Card Company
const (
	Mastercard = "51"
)

// URLs
const (
	TnotifLocal = "http://localhost:8082/sendMail"
)
