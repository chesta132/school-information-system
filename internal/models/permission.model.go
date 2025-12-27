package models

type Permission struct {
	Id
	Name        string             `gorm:"unique" json:"name" example:"role full manage"`
	Resource    PermissionResource `gorm:"type:permission_resource" json:"resource"`
	Description string             `json:"description" example:"Full access to manage role of users"`
	Actions     []PermissionAction `gorm:"type:text;serializer:json" json:"actions"` // []("create", "read", "update", "delete")

	Timestamp
}
