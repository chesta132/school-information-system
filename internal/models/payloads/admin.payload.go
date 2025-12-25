package payloads

import (
	"school-information-system/internal/models"
	"time"
)

type RequestInitiateAdmin struct {
	Key        string    `json:"key" validate:"required"`
	TargetID   string    `json:"target_id" validate:"required"`
	StaffRole  string    `json:"staff_role" validate:"required"`
	EmployeeID string    `json:"employee_id" validate:"required"`
	JoinedAt   time.Time `json:"joined_at" validate:"required"`
}

type RequestSetRole struct {
	TargetID   string          `json:"target_id" validate:"required"`
	TargetRole models.UserRole `json:"target_role" validate:"required,oneof=student teacher admin"`

	StudentData *RequestSetRoleStudent `json:"student_data" validate:"required_if=TargetRole student"`
	TeacherData *RequestSetRoleTeacher `json:"teacher_data" validate:"required_if=TargetRole teacher"`
	AdminData   *RequestSetRoleAdmin   `json:"admin_data" validate:"required_if=TargetRole admin"`
}

type RequestSetRoleStudent struct {
	ClassID   string   `json:"class_id" validate:"required"`
	ParentIDs []string `json:"parent_ids" validate:"required,min=2,max=2"`
	NISN      string   `validate:"required"`
}

type RequestSetRoleTeacher struct {
	SubjectIDs []string  `json:"subject_ids" validate:"required"`
	NUPTK      string    `validate:"required"`
	EmployeeID string    `json:"employee_id" validate:"required"`
	JoinedAt   time.Time `json:"joined_at" validate:"required"`
}

type RequestSetRoleAdmin struct {
	StaffRole  string    `json:"staff_role" validate:"required"`
	EmployeeID string    `json:"employee_id" validate:"required"`
	JoinedAt   time.Time `json:"joined_at" validate:"required"`
}
