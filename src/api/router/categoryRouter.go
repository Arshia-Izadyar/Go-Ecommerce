package router

import (
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/handler"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/gin-gonic/gin"
)

func CategoryRouter(r *gin.RouterGroup, cfg *config.Config) {
	h := handler.NewCategoryHandler(cfg)
	r.POST("/", h.CreateCategory)
}
