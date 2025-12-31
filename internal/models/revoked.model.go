package models

import (
	"time"
)

// do not send revoked data to client
type Revoked struct {
	Id
	Token  string `gorm:"uniqueIndex;<-:create;not null" json:"-"`
	Reason string `gorm:"not null" json:"-"`

	RevokedUntil time.Time `gorm:"index;<-:create;not null" json:"-"`
	Timestamp
}
