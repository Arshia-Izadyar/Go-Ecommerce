package services

import (
	"fmt"
	"time"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/dto"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/constants"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/cache"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/pkg/service_errors"
	"github.com/redis/go-redis/v9"
)

type OtpService struct {
	Redis *redis.Client
	Cfg   *config.Config
}

func NewOtpService(cfg *config.Config) *OtpService {
	redis := cache.GetRedis()
	return &OtpService{
		Redis: redis,
		Cfg:   cfg,
	}
}

func (os *OtpService) SetOtp(req *dto.RequestOtpDTO, otp string) error {
	key := fmt.Sprintf("%s_%s", constants.OTP_REDIS_PREFIX, req.PhoneNumber)

	value := &dto.OtpDTO{
		Value: otp,
		Used:  false,
	}

	v, err := cache.Get[dto.OtpDTO](key)
	if err == nil && v.Used {
		return &service_errors.ServiceError{
			EndUserMessage: service_errors.OtpUsed,
		}
	} else if err == nil && !v.Used {
		return &service_errors.ServiceError{
			EndUserMessage: service_errors.OtpExists,
		}
	}
	err = cache.Set[dto.OtpDTO](key, *value, (os.Cfg.Otp.ExpireTime * time.Minute))
	if err != nil {
		return &service_errors.ServiceError{
			EndUserMessage: service_errors.OtpSetError,
		}
	}
	return nil
}

func (os *OtpService) ValidateOtp(phone_number, otp string) *service_errors.ServiceError {
	key := fmt.Sprintf("%s_%s", constants.OTP_REDIS_PREFIX, phone_number)
	value, err := cache.Get[dto.OtpDTO](key)
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpDoesNotExists}
	}
	if value.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	} else if !value.Used && value.Value != otp {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpInvalid}
	} else if !value.Used && value.Value == otp {
		otpDTO := dto.OtpDTO{
			Value: value.Value,
			Used:  true,
		}
		err = cache.Set[dto.OtpDTO](key, otpDTO, time.Second*2)
		if err != nil {
			return &service_errors.ServiceError{EndUserMessage: service_errors.OtpSetError}
		}
	}
	return nil
}
