package repos

import (
	"context"
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

type Parent struct {
	db *gorm.DB
	create[models.Parent]
	read[models.Parent]
	update[models.Parent]
	delete[models.Parent]
}

func NewParent(db *gorm.DB) *Parent {
	return &Parent{db, create[models.Parent]{db}, read[models.Parent]{db}, update[models.Parent]{db}, delete[models.Parent]{db}}
}

func (r *Parent) WithTx(tx *gorm.DB) *Parent {
	return NewParent(tx)
}

func (r *Parent) DB() *gorm.DB {
	return r.db
}

func (r *Parent) GetStudents(ctx context.Context, parentID string) ([]models.User, error) {
	var students []models.User
	err := r.db.Model(new(models.User)).WithContext(ctx).
		Joins("JOIN students ON students.user_id = users.id").
		Joins("JOIN student_parents sp ON sp.student_id = students.id").
		Where("sp.parent_id = ?", parentID).
		Preload("StudentProfile").
		Preload("StudentProfile.Parents").
		Find(&students).Error

	return students, err
}
