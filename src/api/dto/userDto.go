package dto

type CreateUserDTO struct {
	Username        string `json: "username"`
	Password        string `json: "password"`
	PasswordConfirm string `json: "passwordConfirm"`
	PhoneNumber     string `json: "phoneNumber"`
}

type VerifyUserDTO struct {
	OtpCode     string `json: "otpCode"`
	PhoneNumber string `json: "phoneNumber"`
}

type LoginRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutRequestDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
