package constants

const (
	DEFAULT_ROLE_NAME = "user"
	ADMIN_ROLE_NAME   = "user"
	ADMIN_USERNAME    = "admin"

	// otp
	OTP_REDIS_PREFIX = "otp"

	// jwt
	USER_ID_KEY          = "_userId"
	USERNAME_KEY         = "_username"
	USER_ROLES_KEY       = "_roles"
	TOKEN_TYPE_KEY       = "_type"
	ACCESS_TOKEN_EXPIRE  = "_at_expire"
	REFRESH_TOKEN_EXPIRE = "_rt_expire"
	USER_REFRESH_KEY     = "_refresh_token"

	ACCESS_TOKEN  = "_ac_token"
	REFRESH_TOKEN = "_rf_token"
)
