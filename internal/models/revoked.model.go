package models

import (
	"time"
)

// do not send revoked data to client
type Revoked struct {
	Id
	Token  string `gorm:"uniqueIndex;<-:create" json:"-"`
	Reason string `json:"-"`

	RevokedUntil time.Time `gorm:"index;<-:create" json:"-"`
	Timestamp
}
