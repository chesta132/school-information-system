package repos

import (
	"context"
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

type Class struct {
	db *gorm.DB
	create[models.Class]
	read[models.Class]
	update[models.Class]
	delete[models.Class]
}

func NewClass(db *gorm.DB) *Class {
	return &Class{db, create[models.Class]{db}, read[models.Class]{db}, update[models.Class]{db}, delete[models.Class]{db}}
}

func (r *Class) WithTx(tx *gorm.DB) *Class {
	return NewClass(tx)
}

func (r *Class) DB() *gorm.DB {
	return r.db
}

func (s *Class) GetFormTeacher(ctx context.Context, classID string) (models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).Preload("TeacherProfile").
		Joins("JOIN teachers ON teachers.user_id = users.id").
		Joins("JOIN classes ON classes.form_teacher_id = teachers.id").
		Where("classes.id = ?", classID).
		First(&user).Error
	return user, err
}

func (s *Class) GetStudents(ctx context.Context, classID string) ([]models.User, error) {
	var students []models.User
	err := s.db.WithContext(ctx).Preload("StudentProfile").
		Joins("JOIN students ON students.user_id = users.id").
		Where("students.class_id = ?", classID).
		Find(&students).Error
	return students, err
}
