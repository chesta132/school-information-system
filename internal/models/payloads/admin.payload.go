package payloads

import (
	"time"
)

type RequestInitiateAdmin struct {
	Key        string    `json:"key" validate:"required"`
	TargetID   string    `json:"target_id" validate:"required"`
	StaffRole  string    `json:"staff_role" validate:"required"`
	EmployeeID string    `json:"employee_id" validate:"required"`
	JoinedAt   time.Time `json:"joined_at" validate:"required"`
}
