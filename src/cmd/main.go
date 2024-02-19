package main

import (
	"log"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/cache"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/database"
)

func main() {
	var cfg *config.Config = config.GetConfig()

	err := database.InitDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = cache.InitRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// migrations.Init_01_user_roles()
	api.InitApp(cfg)
}
