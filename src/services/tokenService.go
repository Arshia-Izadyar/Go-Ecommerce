package services

import (
	"fmt"
	"time"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/dto"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/constants"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/cache"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/database"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/models"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/pkg/service_errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type TokenService struct {
	Redis *redis.Client
	cfg   *config.Config
}

func NewTokenService(cfg *config.Config) *TokenService {
	redis := cache.GetRedis()
	return &TokenService{
		Redis: redis,
		cfg:   cfg,
	}
}

func (ts *TokenService) GenerateToken(req *dto.TokenDTO) (*dto.TokenDetailsDTO, error) {
	tokenDetails := dto.TokenDetailsDTO{}
	tokenDetails.AccessTokenExpire = time.Now().Add(time.Minute * ts.cfg.JWT.AccessTokenExpireDuration).Unix()
	tokenDetails.RefreshTokenExpire = time.Now().Add(time.Minute * ts.cfg.JWT.RefreshTokenExpireDuration).Unix()

	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims[constants.USERNAME_KEY] = req.Username
	accessTokenClaims[constants.USER_ROLES_KEY] = req.Roles
	accessTokenClaims[constants.TOKEN_TYPE_KEY] = constants.ACCESS_TOKEN
	accessTokenClaims[constants.ACCESS_TOKEN_EXPIRE] = tokenDetails.AccessTokenExpire

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	var err error
	tokenDetails.AccessToken, err = tk.SignedString([]byte(ts.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims[constants.USERNAME_KEY] = req.Username
	refreshTokenClaims[constants.USER_ID_KEY] = req.UserId
	refreshTokenClaims[constants.USER_ROLES_KEY] = req.Roles
	refreshTokenClaims[constants.TOKEN_TYPE_KEY] = constants.REFRESH_TOKEN
	refreshTokenClaims[constants.REFRESH_TOKEN_EXPIRE] = tokenDetails.RefreshTokenExpire

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	tokenDetails.RefreshToken, err = rt.SignedString([]byte(ts.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}
	return &tokenDetails, nil
}

func (ts *TokenService) ValidateAccessToken(token, tokenType string) (*jwt.Token, error) {
	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
		}
		return []byte(ts.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if tokenType == constants.ACCESS_TOKEN {

		if tokenClaims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid && !ts.IsBlackListed(token) {
			expTime := time.Unix(int64(tokenClaims[constants.ACCESS_TOKEN_EXPIRE].(float64)), 0)
			nowTime := time.Now()

			if nowTime.After(expTime) {
				return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
			}
		}
	} else if tokenType == constants.REFRESH_TOKEN {

		if tokenClaims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid && !ts.IsBlackListed(token) {
			expTime := time.Unix(int64(tokenClaims[constants.REFRESH_TOKEN_EXPIRE].(float64)), 0)
			nowTime := time.Now()

			if nowTime.After(expTime) {
				return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
			}
		}
	}
	return tk, nil
}

func (ts *TokenService) GetClaims(token, tokenType string) (map[string]interface{}, error) {
	claimMap := map[string]interface{}{}
	verification, err := ts.ValidateAccessToken(token, tokenType)
	if err != nil {
		return nil, err
	}
	claims, ok := verification.Claims.(jwt.MapClaims)
	if ok && verification.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessage: service_errors.ClaimNotFound}
}

func (ts *TokenService) AddToBlacklist(token, tokenType string) (bool, error) {

	claims, err := ts.GetClaims(token, tokenType)
	if err != nil {
		return false, err
	}

	if tokenType == constants.ACCESS_TOKEN {
		expTime := time.Unix(int64(claims[constants.ACCESS_TOKEN_EXPIRE].(float64)), 0)
		nowTime := time.Now()
		timeUntilExpire := expTime.Sub(nowTime)
		fmt.Println(timeUntilExpire)

		err := cache.Set(token, true, timeUntilExpire)
		if err != nil {
			return false, err
		}
		return true, nil

	} else if tokenType == constants.REFRESH_TOKEN {
		expTime := time.Unix(int64(claims[constants.REFRESH_TOKEN_EXPIRE].(float64)), 0)
		nowTime := time.Now()
		timeUntilExpire := expTime.Sub(nowTime)
		fmt.Println(timeUntilExpire)
		err := cache.Set(token, true, timeUntilExpire)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}

func (ts *TokenService) IsBlackListed(token string) bool {
	_, err := cache.Get[bool](token)
	return err == nil
}

func (ts *TokenService) RefreshToken(req *dto.RefreshTokenDTO) (*dto.TokenDetailsDTO, error) {
	tk, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {

		return []byte(ts.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok || !tk.Valid {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
	}
	if ts.IsBlackListed(req.RefreshToken) || claims[constants.TOKEN_TYPE_KEY] != constants.REFRESH_TOKEN {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
	}

	expTime := time.Unix(int64(claims[constants.REFRESH_TOKEN_EXPIRE].(float64)), 0)
	if time.Now().After(expTime) {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
	}

	ts.AddToBlacklist(req.RefreshToken, constants.REFRESH_TOKEN)

	userId := claims[constants.USER_ID_KEY]
	var user *models.User
	db := database.GetDB()
	err = db.Model(&models.User{}).Where("id = ?", userId).Preload("UserRoles.Role").Find(&user).Error
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UserNotFound}
	}
	tokenDTO := &dto.TokenDTO{
		Username: user.UserName,
		UserId:   user.Id,
	}
	if len(user.UserRoles) >= 1 {
		for _, role := range user.UserRoles {
			tokenDTO.Roles = append(tokenDTO.Roles, role.Role.Name)
		}
	}
	key := fmt.Sprintf("%s%s", user.UserName, constants.USER_REFRESH_KEY)

	token, err := ts.GenerateToken(tokenDTO)
	if err != nil {
		return nil, err
	}
	err = cache.Set[string](key, token.RefreshToken, ts.cfg.JWT.RefreshTokenExpireDuration*time.Minute)
	if err != nil {
		return nil, err
	}
	return token, nil
}
