package router

import (
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/handler"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup, cfg *config.Config) {
	UsrHandler := handler.NewUserHandler(cfg)
	r.POST("/register", UsrHandler.RegisterUser)
	r.POST("/verify", UsrHandler.VerifyUser)
	r.POST("/login", UsrHandler.LoginUser)
	r.POST("/logout", UsrHandler.LogoutUser)
}
