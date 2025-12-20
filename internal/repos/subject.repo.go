package repos

import (
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

type Subject struct {
	db *gorm.DB
	create[models.Subject]
	read[models.Subject]
	update[models.Subject]
	delete[models.Subject]
}

func NewSubject(db *gorm.DB) *Subject {
	return &Subject{db, create[models.Subject]{db}, read[models.Subject]{db}, update[models.Subject]{db}, delete[models.Subject]{db}}
}

func (r *Subject) WithTx(tx *gorm.DB) *Subject {
	return NewSubject(tx)
}

func (r *Subject) DB() *gorm.DB {
	return r.db
}
