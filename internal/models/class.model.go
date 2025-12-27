package models

type Class struct {
	Id
	Name          string   `json:"name" example:"10 TJKT 3"`
	FormTeacherID string   `json:"-"`
	FormTeacher   *Teacher `json:"form_teacher,omitempty" swaggerignore:"true"`

	Timestamp
}
