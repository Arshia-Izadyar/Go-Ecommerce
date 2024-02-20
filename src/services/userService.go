package services

import (
	"fmt"
	"time"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/dto"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/common"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/constants"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/cache"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/database"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/models"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/pkg/service_errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB    *gorm.DB
	Otp   *OtpService
	Token *TokenService
	Cfg   *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	db := database.GetDB()
	otp := NewOtpService(cfg)
	token := NewTokenService(cfg)
	return &UserService{
		DB:    db,
		Otp:   otp,
		Token: token,
		Cfg:   cfg,
	}
}

func (us *UserService) CheckByUsername(username string) bool {
	var exists bool
	us.DB.Model(&models.User{}).Select("count(*) > 0").Where("user_name = ?", username).Find(&exists)
	return exists
}

func (us *UserService) CheckByPhone(phone string) bool {
	var exists bool
	us.DB.Model(&models.User{}).Select("count(*) > 0").Where("phone_number = ?", phone).Find(&exists)
	return exists
}

func (us *UserService) GetDefaultRole() (*models.Role, error) {
	var role *models.Role
	err := us.DB.Model(&models.Role{}).Where("name = ?", constants.DEFAULT_ROLE_NAME).Find(&role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (us *UserService) RegisterUser(req *dto.CreateUserDTO) error {
	phoneExists := us.CheckByPhone(req.PhoneNumber)
	usernameExists := us.CheckByUsername(req.Username)
	if usernameExists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UsernameExists}
	} else if phoneExists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PhoneNumberExists}
	}
	if req.Password != req.PasswordConfirm {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PasswordsConfirmWrong}
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.CantCreateUser}
	}

	role, err := us.GetDefaultRole()
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.CantCreateUser}
	}

	user := &models.User{
		UserName:    req.Username,
		PhoneNumber: req.PhoneNumber,
		Password:    string(bs),
	}

	tx := us.DB.Begin()

	err = tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: service_errors.CantCreateUser}
	}
	userRole := &models.UserRole{UserId: user.Id, RoleId: role.Id}
	err = tx.Create(userRole).Error
	if err != nil {
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: service_errors.CantCreateUser}
	}

	otp := common.GenerateOtp()
	otpDTO := &dto.RequestOtpDTO{
		PhoneNumber: user.PhoneNumber,
	}
	err = us.Otp.SetOtp(otpDTO, otp)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: "otp"}
	}

	tx.Commit()

	return nil
}

func (us *UserService) UserVerify(req *dto.VerifyUserDTO) (bool, error) {

	var user models.User
	err := us.DB.Model(&models.User{}).Where("phone_number = ?", req.PhoneNumber).First(&user).Error
	if err != nil {
		return false, err
	}

	if user.Verified {
		return false, &service_errors.ServiceError{EndUserMessage: service_errors.UserAlreadyVerified}
	}

	key := fmt.Sprintf("%s_%s", constants.OTP_REDIS_PREFIX, req.PhoneNumber)

	redisOtp, err := cache.Get[dto.OtpDTO](key)
	if err != nil {
		return false, &service_errors.ServiceError{EndUserMessage: service_errors.OtpInvalid}
	}

	if redisOtp.Value != req.OtpCode {
		return false, &service_errors.ServiceError{EndUserMessage: service_errors.OtpInvalid}
	} else if redisOtp.Used {
		return false, &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	}
	///
	user.Verified = true
	us.DB.Save(user)

	err = cache.Set(key, "", time.Second)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (us *UserService) LoginUser(req *dto.LoginRequestDto) (*dto.TokenDetailsDTO, error) {
	var user models.User
	err := us.DB.Model(&models.User{}).Where("user_name = ?", req.Username).Preload("UserRoles.Role").First(&user).Error
	if err != nil {
		return nil, err
	}

	if !user.Verified {
		return nil, &service_errors.ServiceError{EndUserMessage: "user is not verified please verify first"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.WrongPassword}
	}

	key := fmt.Sprintf("%s%s", req.Username, constants.USER_REFRESH_KEY)
	ext, err := cache.Get[string](key)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if ext != nil {
		existingToken := *ext
		tokenDTO := &dto.RefreshTokenDTO{
			RefreshToken: existingToken,
		}
		claims, err := us.Token.RefreshToken(tokenDTO)
		if err == nil {
			return claims, nil
		}
	}

	tokenDto := &dto.TokenDTO{
		Username: user.UserName,
		UserId:   user.Id,
	}
	fmt.Println(user.UserRoles)

	if len(user.UserRoles) > 0 {
		for _, role := range user.UserRoles {
			fmt.Println(role.Role.Name)
			fmt.Println(role.Role.Id)

			tokenDto.Roles = append(tokenDto.Roles, role.Role.Name)
		}
	}

	token, err := us.Token.GenerateToken(tokenDto)
	if err != nil {
		return nil, err
	}
	err = cache.Set[string](key, token.RefreshToken, us.Cfg.JWT.RefreshTokenExpireDuration*time.Minute)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (us *UserService) Logout(req *dto.LogoutRequestDto) error {
	_, err := us.Token.AddToBlacklist(req.AccessToken, constants.ACCESS_TOKEN)
	if err != nil {
		return err
	}
	_, err = us.Token.AddToBlacklist(req.RefreshToken, constants.REFRESH_TOKEN)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) RefreshToken(req *dto.RefreshTokenDTO) (*dto.TokenDetailsDTO, error) {
	token, err := us.Token.RefreshToken(req)
	if err != nil {
		return nil, err
	}
	return token, nil
}
