package models

import (
	"time"

	"gorm.io/gorm"
)

type Id struct {
	ID string `gorm:"default:gen_random_uuid()" json:"id" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
}

type Timestamp struct {
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitzero" example:"2006-01-02T15:04:05Z07:00"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitzero" example:"2006-01-02T15:04:05Z07:00"`
}

type TimestampJoinTime struct {
	JoinedAt  time.Time `gorm:"not null" json:"joined_at" example:"2006-01-02T15:04:05Z07:00"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitzero" example:"2006-01-02T15:04:05Z07:00"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitzero" example:"2006-01-02T15:04:05Z07:00"`
}

type TimestampArchivable struct {
	DeletedAt gorm.DeletedAt `gorm:"index" json:"archived_at,omitzero" swaggertype:"string" example:"2006-01-02T15:04:05Z07:00"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at,omitzero" example:"2006-01-02T15:04:05Z07:00"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitzero" example:"2006-01-02T15:04:05Z07:00"`
}
