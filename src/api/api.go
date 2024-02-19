package api

import (
	"fmt"
	"log"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/router"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/gin-gonic/gin"
)

func InitApp(cfg *config.Config) {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	registerRoutes(r, cfg)

	err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatal(err)
	}

}

func registerRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/users")
	router.UserRouter(users, cfg)

}
