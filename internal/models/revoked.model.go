package models

import (
	"time"
)

// do not send revoked data to client
type Revoked struct {
	ID     string `gorm:"default:gen_random_uuid()" json:"-"`
	Token  string `gorm:"uniqueIndex;<-:create" json:"-"`
	Reason string `json:"-"`

	RevokedUntil time.Time `gorm:"index;<-:create" json:"-"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"-"`
}
