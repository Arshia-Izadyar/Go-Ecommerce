package migrations

import "github.com/Arshia-Izadyar/Go-Ecommerce/src/data/database"

func user_product_02_migration() {
	db := database.GetDB()
	createTables(db)
}
