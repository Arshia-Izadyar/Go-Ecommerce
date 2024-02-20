package handler

import (
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/dto"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/helper"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	s := services.NewUserService(cfg)
	return &UserHandler{
		service: s,
	}
}

func (uh *UserHandler) RegisterUser(ctx *gin.Context) {
	req := &dto.CreateUserDTO{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithValidationError(err, 400, false))
		return
	}
	err = uh.service.RegisterUser(req)
	if err != nil {

		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithError(err, 400, false))
		return
	}
	ctx.JSON(200, helper.GenerateResponse("otp code sent !", 200, true))

}

func (uh *UserHandler) VerifyUser(ctx *gin.Context) {
	req := &dto.VerifyUserDTO{}
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithValidationError(err, 400, false))
		return
	}
	_, err = uh.service.UserVerify(req)

	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithError(err, 400, false))
		return
	}
	ctx.JSON(200, helper.GenerateResponse("user verified", 200, true))
}

func (uh *UserHandler) LoginUser(ctx *gin.Context) {
	req := &dto.LoginRequestDto{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithValidationError(err, 400, false))
		return
	}
	token, err := uh.service.LoginUser(req)
	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithError(err, 400, false))
		return
	}
	ctx.JSON(200, helper.GenerateResponse(token, 200, true))

}

// FIXME: add refreshToken from cookie and access token from header
func (uh *UserHandler) LogoutUser(ctx *gin.Context) {
	req := &dto.LogoutRequestDto{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithValidationError(err, 400, false))
		return
	}
	err = uh.service.Logout(req)
	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithError(err, 400, false))
		return
	}
	ctx.JSON(204, helper.GenerateResponse("successfully logged out", 204, true))

}

func (uh *UserHandler) RefreshToken(ctx *gin.Context) {
	req := &dto.RefreshTokenDTO{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithValidationError(err, 400, false))
		return
	}
	token, err := uh.service.RefreshToken(req)
	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithError(err, 400, false))
		return
	}
	ctx.JSON(200, helper.GenerateResponse(token, 200, true))

}
