package database

import (
	"log"
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := CreateEnum(db, "user_role", []string{
		string(models.RoleStudent),
		string(models.RoleTeacher),
		string(models.RoleAdmin),
		string(models.RoleUnsetted)},
	); err != nil {
		log.Fatal("[MIGRATE] failed to create user role enum", err.Error())
	}

	if err := CreateEnum(db, "user_gender", []string{
		string(models.GenderMale),
		string(models.GenderFemale),
	}); err != nil {
		log.Fatal("[MIGRATE] failed to create user gender enum", err.Error())
	}

	if err := CreateEnum(db, "permission_resource", []string{
		string(models.ResourceRole),
		string(models.ResourcePermission),

		string(models.ResourceAdmin),
		string(models.ResourceTeacher),
		string(models.ResourceStudent),

		string(models.ResourceSubject),
		string(models.ResourceClass),
	}); err != nil {
		log.Fatal("[MIGRATE] failed to create permission resource enum", err.Error())
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
		log.Fatal("[MIGRATE] failed to migrate databases model", err.Error())
	}
}
