package models

type Permission struct {
	Id
	Name        string             `gorm:"unique" json:"name"`
	Resource    PermissionResource `gorm:"type:permission_resource" json:"resource"`
	Description string             `json:"description"`
	Actions     []PermissionAction `gorm:"type:text;serializer:json" json:"actions"` // []("create", "read", "update", "delete")

	Timestamp
}
