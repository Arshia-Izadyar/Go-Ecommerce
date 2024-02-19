package dto

type TokenDTO struct {
	Username string   `json:"username"`
	UserId   int      `json:"userId"`
	Roles    []string `json:"roles"`
}

type TokenDetailsDTO struct {
	AccessToken        string `json: "accessToken"`
	RefreshToken       string `json: "refreshToken"`
	AccessTokenExpire  int64  `json: "accessTokenExpire"`
	RefreshTokenExpire int64  `json: "refreshTokenExpire"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refreshToken"`
}
