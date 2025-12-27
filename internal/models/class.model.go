package models

type Class struct {
	Id
	Name          string   `json:"name"`
	FormTeacherID string   `json:"-"`
	FormTeacher   *Teacher `json:"form_teacher,omitempty"`

	Timestamp
}
