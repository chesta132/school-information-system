package payloads

import (
	"school-information-system/internal/models"
	"time"
)

type RequestInitiateAdmin struct {
	Key        string    `json:"key" validate:"required" example:"super secret"`
	TargetID   string    `json:"target_id" validate:"required,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	StaffRole  string    `json:"staff_role" validate:"required" example:"developer"`
	EmployeeID string    `json:"employee_id" validate:"required" example:"DEV002"`
	JoinedAt   time.Time `json:"joined_at" validate:"required" example:"2006-01-02T15:04:05Z07:00"`
}

type RequestSetRole struct {
	TargetID   string          `json:"target_id" validate:"required,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	TargetRole models.UserRole `json:"target_role" validate:"required,user_role"`

	// empty ig target_role not student
	StudentData *RequestSetRoleStudent `json:"student_data" prefix:"student_data." validate:"required_if=TargetRole student"`
	// empty ig target_role not teacher
	TeacherData *RequestSetRoleTeacher `json:"teacher_data" prefix:"teacher_data." validate:"required_if=TargetRole teacher"`
	// empty ig target_role not admin
	AdminData *RequestSetRoleAdmin `json:"admin_data" prefix:"admin_data." validate:"required_if=TargetRole admin"`
}

type RequestSetRoleStudent struct {
	ClassID   string   `json:"class_id" validate:"required,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	ParentIDs []string `json:"parent_ids" validate:"required,min=2,max=2,dive,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6,479b5b5f-81b1-4669-91a5-b5bf69e597c7"`
	NISN      string   `validate:"required" example:"0091913711"`
}

type RequestSetRoleTeacher struct {
	SubjectIDs []string  `json:"subject_ids" validate:"required,min=1,max=1,dive,uuid4" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6,479b5b5f-81b1-4669-91a5-b5bf69e597c7"`
	NUPTK      string    `validate:"required" example:"1234567890123456"`
	EmployeeID string    `json:"employee_id" validate:"required" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	JoinedAt   time.Time `json:"joined_at" validate:"required" example:"2006-01-02T15:04:05Z07:00"`
}

type RequestSetRoleAdmin struct {
	StaffRole  string    `json:"staff_role" validate:"required" example:"developer"`
	EmployeeID string    `json:"employee_id" validate:"required" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	JoinedAt   time.Time `json:"joined_at" validate:"required" example:"2006-01-02T15:04:05Z07:00"`
}
