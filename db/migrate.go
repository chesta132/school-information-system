package db

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("school-app: failed to migrate databases model", err.Error())
	}

	if err := CreateEnum(db, "user_role", []string{"student", "teacher", "admin"}); err != nil {
		log.Fatal("school-app: failed to create user role enum", err.Error())
	}

	if err := CreateEnum(db, "user_gender", []string{"male", "female"}); err != nil {
		log.Fatal("school-app: failed to create user gender enum", err.Error())
	}

	if err := CreateEnum(db, "permission_action", []string{"create", "read", "update", "delete"}); err != nil {
		log.Fatal("school-app: failed to create permission action enum", err.Error())
	}
}
