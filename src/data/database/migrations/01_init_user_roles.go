package migrations

import (
	"log"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/constants"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/database"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Init_01_user_roles() {
	db := database.GetDB()
	createTables(db)
}

func createTables(db *gorm.DB) {
	tables := []interface{}{}

	tables = addTableIfNotExist(models.User{}, db, tables)
	tables = addTableIfNotExist(models.Role{}, db, tables)
	tables = addTableIfNotExist(models.UserRole{}, db, tables)
	tables = addTableIfNotExist(models.Category{}, db, tables)
	tables = addTableIfNotExist(models.Product{}, db, tables)
	tables = addTableIfNotExist(models.ProductCategory{}, db, tables)
	tables = addTableIfNotExist(models.UserWishList{}, db, tables)

	err := db.Migrator().CreateTable(tables...)
	if err != nil {
		db.Rollback()
		log.Fatal(err)
	}
	createInitData(db)

}

func addTableIfNotExist(model interface{}, db *gorm.DB, tables []interface{}) []interface{} {

	if !db.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func createRoleIfNotExists(db *gorm.DB, r *models.Role) int {
	var exists bool
	var role models.Role
	db.Model(models.Role{}).Where("name = ?", r.Name).Select("COUNT(*) > 0").Find(&exists)
	if !exists {
		db.Model(models.Role{}).Create(&models.Role{Name: r.Name}).Where("name = ?", r.Name).First(&role)
		return role.Id
	}
	return 0
}

func createInitData(db *gorm.DB) {
	defaultRole := &models.Role{
		Name: constants.DEFAULT_ROLE_NAME,
	}
	createRoleIfNotExists(db, defaultRole)
	adminRole := &models.Role{
		Name: constants.ADMIN_ROLE_NAME,
	}
	adminId := createRoleIfNotExists(db, adminRole)

	admin := &models.User{
		UserName:    constants.ADMIN_USERNAME,
		PhoneNumber: "+989108624707",
		Verified:    true,
	}
	bs, err := bcrypt.GenerateFromPassword([]byte("a123"), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	admin.Password = string(bs)
	createAdmin(admin, db, adminId)

}

func createAdmin(model *models.User, db *gorm.DB, roleId int) {
	var exists bool
	db.Model(&models.User{}).Where("user_name = ? ", model.UserName).Select("COUNT(*) > 0").Find(&exists)
	if !exists {
		db.Model(&models.User{}).Create(model)
		userRole := &models.UserRole{
			UserId: model.Id,
			RoleId: roleId,
		}
		db.Create(&userRole)
	}
}
