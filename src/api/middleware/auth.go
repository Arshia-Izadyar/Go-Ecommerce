package middleware

import (
	"fmt"
	"strings"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/helper"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/constants"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/pkg/service_errors"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/services"
	"github.com/gin-gonic/gin"
)

func Authenticate(cfg *config.Config) gin.HandlerFunc {
	tokenService := services.NewTokenService(cfg)
	var err error
	var tokenClaims map[string]interface{}
	return func(ctx *gin.Context) {
		header := ctx.GetHeader(constants.AUTH_HEADER)
		if header == "" {
			ctx.AbortWithStatusJSON(401, helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "no token provided"}, 401, false))
			return
		}
		token := strings.Split(header, " ")[1]
		if token == "" || !strings.HasPrefix(header, "Bearer") {
			ctx.AbortWithStatusJSON(401, helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: "invalid token provided"}, 401, false))
			return
		}
		tokenClaims, err = tokenService.GetClaims(token, constants.ACCESS_TOKEN)
		if err != nil {
			ctx.AbortWithStatusJSON(401, helper.GenerateResponseWithError(&service_errors.ServiceError{EndUserMessage: err.Error()}, 401, false))
			return
		}
		fmt.Println(tokenClaims)

		ctx.Set(constants.USERNAME_KEY, tokenClaims[constants.USERNAME_KEY])
		ctx.Set(constants.USER_ROLES_KEY, tokenClaims[constants.USER_ROLES_KEY])

		ctx.Next()

	}
}
