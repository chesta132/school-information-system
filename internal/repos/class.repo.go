package repos

import (
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
