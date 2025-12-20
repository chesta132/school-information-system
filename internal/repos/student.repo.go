package repos

import (
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

type Student struct {
	db *gorm.DB
	create[models.Student]
	read[models.Student]
	update[models.Student]
	delete[models.Student]
}

func NewStudent(db *gorm.DB) *Student {
	return &Student{db, create[models.Student]{db}, read[models.Student]{db}, update[models.Student]{db}, delete[models.Student]{db}}
}

func (r *Student) WithTx(tx *gorm.DB) *Student {
	return NewStudent(tx)
}

func (r *Student) DB() *gorm.DB {
	return r.db
}
