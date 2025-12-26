package payloads

import "school-information-system/internal/models"

type RequestGrantPermission struct {
	TargetID     string `json:"target_id" validate:"required"`
	PermissionID string `json:"permission_id" validate:"required"`
}

type RequestRevokePermission struct {
	TargetID     string `json:"target_id" validate:"required"`
	PermissionID string `json:"permission_id" validate:"required"`
}

type RequestCreatePermission struct {
	Name        string                    `json:"name" validate:"required,min=10"`
	Description string                    `json:"description" validate:"required,min=10"`
	Resource    models.PermissionResource `json:"resource" validate:"required,oneof=role permission"`
	Actions     []models.PermissionAction `json:"actions" validate:"required,min=1,dive,oneof=create read update delete"`
}
