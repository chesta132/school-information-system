package models

import "fmt"

type Class struct {
	Id
	Grade         int      `json:"grade" example:"10"`
	Major         string   `json:"major" example:"TJKT"`
	ClassNumber   int      `json:"class_number" example:"3"`
	Name          string   `json:"name" gorm:"-" example:"10 TJKT 3"`
	FormTeacherID *string  `json:"form_teacher_id" gorm:"unique" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	FormTeacher   *Teacher `json:"form_teacher,omitempty" swaggerignore:"true"`

	Timestamp
}

func (c *Class) GetName() string {
	c.Name = fmt.Sprintf("%d %s %d", c.Grade, c.Major, c.ClassNumber)
	return c.Name
}
