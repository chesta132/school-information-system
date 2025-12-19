package models

import "time"

type Subject struct {
	ID   string `gorm:"default:gen_random_uuid()" json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitzero"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitzero"`
}
