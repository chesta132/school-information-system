package repos

import (
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

type Admin struct {
	db *gorm.DB
	create[models.Admin]
	read[models.Admin]
	update[models.Admin]
	delete[models.Admin]
}

func NewAdmin(db *gorm.DB) *Admin {
	return &Admin{db, create[models.Admin]{db}, read[models.Admin]{db}, update[models.Admin]{db}, delete[models.Admin]{db}}
}

func (r *Admin) WithTx(tx *gorm.DB) *Admin {
	return NewAdmin(tx)
}

func (r *Admin) DB() *gorm.DB {
	return r.db
}
