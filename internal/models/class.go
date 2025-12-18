package models

import "time"

type Class struct {
	ID            string   `gorm:"default:gen_random_uuid()" json:"id"`
	Name          string   `json:"name"`
	FormTeacherID string   `json:"-"`
	FormTeacher   *Teacher `json:"form_teacher,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
