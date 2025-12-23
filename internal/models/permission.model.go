package models

import "time"

type Permission struct {
	ID          string             `gorm:"default:gen_random_uuid()" json:"id"`
	Name        string             `gorm:"unique" json:"name"`
	Resource    PermissionResource `gorm:"type:permission_resource" json:"resource"`
	Description string             `json:"description"`
	Actions     []PermissionAction `gorm:"type:permission_action[]" json:"actions"` // []("create", "read", "update", "delete")
	AuthorID    string             `json:"-"`
	Author      *Admin             `json:"author,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitzero"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitzero"`
}
