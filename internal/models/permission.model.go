package models

type Permission struct {
	Id
	Name        string             `gorm:"unique;not null" json:"name" example:"role full manage"`
	Resource    PermissionResource `gorm:"type:permission_resource;not null" json:"resource"`
	Description string             `gorm:"not null" json:"description" example:"Full access to manage role of users"`
	Actions     []PermissionAction `gorm:"type:text;serializer:json;not null" json:"actions"` // []("create", "read", "update", "delete")

	Timestamp
}
