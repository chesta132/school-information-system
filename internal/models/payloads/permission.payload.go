package payloads

type RequestGrantPermission struct {
	TargetID     string `json:"target_id" validate:"required"`
	PermissionID string `json:"permission_id" validate:"required"`
}

type RequestRevokePermission struct {
	TargetID     string `json:"target_id" validate:"required"`
	PermissionID string `json:"permission_id" validate:"required"`
}
