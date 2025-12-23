package database

import (
	"log"
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := CreateEnum(db, "user_role", []string{"student", "teacher", "admin", "unsetted"}); err != nil {
		log.Fatal("school-app: failed to create user role enum", err.Error())
	}

	if err := CreateEnum(db, "user_gender", []string{"male", "female"}); err != nil {
		log.Fatal("school-app: failed to create user gender enum", err.Error())
	}

	if err := CreateEnum(db, "permission_action", []string{"create", "read", "update", "delete"}); err != nil {
		log.Fatal("school-app: failed to create permission action enum", err.Error())
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Teacher{},
		&models.Admin{},
		&models.Teacher{},
		&models.Permission{},
		&models.Class{},
		&models.Subject{},
		&models.Revoked{},
	); err != nil {
		log.Fatal("school-app: failed to migrate databases model", err.Error())
	}
}
