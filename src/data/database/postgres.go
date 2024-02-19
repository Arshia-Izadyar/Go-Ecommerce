package database

import (
	"fmt"
	"log"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error
var SqlDB *gorm.DB

func InitDB(conf *config.Config) error {
	cnn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.DbName,
	)
	SqlDB, err = gorm.Open(postgres.Open(cnn), &gorm.Config{})
	if err != nil {
		return err
	}

	db, err := SqlDB.DB()
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(conf.Postgres.MaxIdleConns)
	db.SetMaxOpenConns(conf.Postgres.MaxOpenConns)
	db.SetConnMaxLifetime(conf.Postgres.ConnMaxLifetime)
	return nil
}

func GetDB() *gorm.DB {
	return SqlDB
}

func CloseDB() {
	db, err := SqlDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
