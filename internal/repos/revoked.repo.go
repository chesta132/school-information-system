package repos

import (
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

type Revoked struct {
	db *gorm.DB
	create[models.Revoked]
	read[models.Revoked]
	update[models.Revoked]
	delete[models.Revoked]
}

func NewRevoked(db *gorm.DB) *Revoked {
	return &Revoked{db, create[models.Revoked]{db}, read[models.Revoked]{db}, update[models.Revoked]{db}, delete[models.Revoked]{db}}
}

func (r *Revoked) WithTx(tx *gorm.DB) *Revoked {
	return NewRevoked(tx)
}

func (r *Revoked) DB() *gorm.DB {
	return r.db
}
