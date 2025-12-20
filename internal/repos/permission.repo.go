package repos

import (
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

type Permission struct {
	db *gorm.DB
	create[models.Permission]
	read[models.Permission]
	update[models.Permission]
	delete[models.Permission]
}

func NewPermission(db *gorm.DB) *Permission {
	return &Permission{db, create[models.Permission]{db}, read[models.Permission]{db}, update[models.Permission]{db}, delete[models.Permission]{db}}
}

func (r *Permission) WithTx(tx *gorm.DB) *Permission {
	return NewPermission(tx)
}

func (r *Permission) DB() *gorm.DB {
	return r.db
}
