package service_errors

const (
	OtpExists        = "otp exits"
	OtpDoesNotExists = "otp doesn't exits"
	OtpUsed          = "otp used"
	OtpInvalid       = "otp invalid"
	OtpSetError      = "can't sent otp right now"

	ClaimNotFound = "claim not found"

	// user
	EmailExists           = "email already exits"
	UsernameExists        = "Username already exits"
	PhoneNumberExists     = "phone number already exits"
	WrongPassword         = "WrongPassword"
	PasswordsConfirmWrong = "password and password confirm don't match"
	UserNotFound          = "User Not Found"
	CantCreateUser        = "cant create user right now"
	UserAlreadyVerified   = "the user is already verified "

	TokenNotPresent = "no token provided"
	TokenExpired    = "token is expired !"
	TokenInvalid    = "provided token is invalid1"
	NotRefreshToken = "provided token is not a refresh token"
	InternalError   = "some thing happened"

	PermissionDenied = "Permission Denied"
)
