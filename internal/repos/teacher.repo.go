package repos

import (
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

type Teacher struct {
	db *gorm.DB
	create[models.Teacher]
	read[models.Teacher]
	update[models.Teacher]
	delete[models.Teacher]
}

func NewTeacher(db *gorm.DB) *Teacher {
	return &Teacher{db, create[models.Teacher]{db}, read[models.Teacher]{db}, update[models.Teacher]{db}, delete[models.Teacher]{db}}
}

func (r *Teacher) WithTx(tx *gorm.DB) *Teacher {
	return NewTeacher(tx)
}

func (r *Teacher) DB() *gorm.DB {
	return r.db
}
