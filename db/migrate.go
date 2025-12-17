package db

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("school-app: failed to migrate databases model", err.Error())
	}
}
