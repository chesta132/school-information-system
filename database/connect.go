package database

import (
	"fmt"
	"log"
	"school-information-system/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DB_HOST, config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_PORT)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Fatal("school-app: failed to connect database", err.Error())
	}

	Migrate(db)

	return db
}
