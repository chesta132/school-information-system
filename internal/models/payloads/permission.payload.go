package payloads

import "school-information-system/internal/models"

type RequestGrantPermission struct {
	TargetID     string `json:"target_id" validate:"required,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	PermissionID string `json:"permission_id" validate:"required,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
}

type RequestRevokePermission struct {
	TargetID     string `json:"target_id" validate:"required,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	PermissionID string `json:"permission_id" validate:"required,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
}

type RequestCreatePermission struct {
	Name        string                    `json:"name" validate:"required,min=10" example:"role setter"`
	Description string                    `json:"description" validate:"required,min=10" example:"permission to set and read role"`
	Resource    models.PermissionResource `json:"resource" validate:"required,permission_resource" example:"role"`
	Actions     []models.PermissionAction `json:"actions" validate:"required,min=1,dive,permission_action" example:"update,read"`
}

type RequestGetPermission struct {
	ID string `uri:"id" validate:"required,uuid4"`
}

type RequestGetPermissions struct {
	Offset    int                         `form:"offset" example:"10"`
	Query     string                      `form:"q" example:"rol"`
	Resources []models.PermissionResource `form:"resource" validate:"dive,permission_resource"`
	Actions   []models.PermissionAction   `form:"action" validate:"dive,permission_action"`
}

type RequestUpdatePermission struct {
	ID          string `validate:"required,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	Name        string `json:"name" validate:"required_without=Description,omitempty,min=10" example:"updated name"`
	Description string `json:"description" validate:"required_without=Name,omitempty,min=10" example:"updated desc"`
}

type RequestDeletePermission struct {
	ID string `uri:"id" validate:"required,uuid4"`
}
