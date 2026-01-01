package payloads

import "school-information-system/internal/models"

type RequestCreateParent struct {
	FullName string            `json:"full_name" validate:"required" example:"Chesta Ardiona"`
	Phone    string            `json:"phone" validate:"required" example:"+6281234567890"`
	Email    string            `json:"email" validate:"required,email" example:"chestaardi4@gmail.com"`
	Gender   models.UserGender `json:"gender" validate:"required,user_gender"`
}
