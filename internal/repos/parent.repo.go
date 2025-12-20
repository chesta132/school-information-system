package repos

import (
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
