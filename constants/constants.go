package constants

// Transaction types
const (
	Deposit = "DEPOSIT_TO_MAIN_ACCOUNT"
)

// Response status
const (
	DepositSuccess = "Deposit completed successfully"
	DepositFailed  = "Failed to make deposit"
	Success        = "Success"
	Verified       = "verified"
)

// Error messages
const (
	PageNotFound                    = "Seems like %v is not available or does not exist"
	RequestHasInvalidFields         = "Mandatory field(s) is missing from request, or is invalid. Please ensure all fields are properly filled."
	CannotReadRequest               = "Cannot read request body"
	CannotParseRequest              = "Unable to parse request"
	UnauthorizedDepositRequest      = "Failed to recognize deposit source. Please use authorized banks to deposit."
	CannotEncodeResponse            = "Failed to encode response."
	CannotEncodeTnotifResponse      = "Bad Response from Email service"
	InsertFailed                    = "Failed to insert to database."
	InitLogFailed                   = "Failed to insert log of this transaction."
	InvalidPassword                 = "Password Min 6 Character"
	InsertTokenFailed               = "Email Token Failed"
	CreateAccountFailed             = "Account Failed"
	CannotParseTnotifRequest        = "Unable to Parse JSON request to Email service"
	NotUniquePhoneNumberAndEmail    = "Unable to Register, Your Phone Number Or Email Has Been Used"
	FailedToFetchHistoryTransaction = "Cannot get history transaction"
	UnauthorizedUser                = "Wrong Email or Password"
	ErrorWhenCallingTnotif          = "Can't Send Email"
	InvalidUserEmail                = "Invalid email"
	UpdateUserDataFailed            = "Error updating customer data"
	RetrieveUserFileFailed          = "Error Retrieving the File"
	VerifyEmailTokenFailed          = "Unable to verify email token: "
	DeleteEmailTokenFailed          = "Unable to delete verified email: "
)

// Hosts
const (
	TnotifLocalhost = "http://localhost:8082"
)

// Endpoints
const (
	TnotifEndpoint                   = "/sendMail"
	HomeEndpoint                     = "/"
	RegisterEndpoint                 = "/register"
	LoginEndpoint                    = "/login"
	CreateVirtualAccountEndpoint     = "/virtualaccount/create"
	EditVirtualAccountEndpoint       = "/virtualaccount/edit"
	TransferVacToMainAccountEndpoint = "/vac/to_main"
	ListAllVacEndpoint               = "/vac/list"
	DeleteVacEndpoint                = "/vac/delete-vac"
	AddVacBalanceEndpoint            = "/vac/add_balance_vac"
	GetTransactionHistoryEndpoint    = "/transaction/history"
	VerifyEmailEndpoint              = "/email/verify-email-token"
	GetUserProfileEndpoint           = "/customers/getprofile"
	UpdateUserProfileEndpoint        = "/customers/updateprofile"
	UpdateUserPhotoEndpoint          = "/customers/updatephoto"
	DepositEndpoint                  = "/deposit"
)

// Headers
const (
	ApplicationJson = "application/json"
)

// Formats
const (
	DateFormat             = "060102"
	UserPhotoFileName      = "myPhoto"
	UserPhotoFolderName    = "temp-images"
	UserPhotoFileExtension = ".png"
	EmailRegex             = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

// Page contents
const (
	HomePage = "Welcome to Tsaving"
)
