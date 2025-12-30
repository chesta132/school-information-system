package payloads

import (
	"school-information-system/internal/models"
)

type RequestSignUp struct {
	FullName   string            `json:"full_name" validate:"required" example:"Chesta Ardiona"`
	Email      string            `json:"email" validate:"required,email" example:"chestaardi4@gmail.com"`
	Password   string            `json:"password" validate:"required,min=8" example:"super.secret871798"`
	Gender     models.UserGender `json:"gender" validate:"required,user_gender"`
	Phone      string            `json:"phone" validate:"required" example:"+6281234567890"`
	RememberMe bool              `json:"remember_me"`
}

type RequestSignIn struct {
	Email      string `json:"email" validate:"required,email" example:"chestaardi4@gmail.com"`
	Password   string `json:"password" validate:"required" example:"super.secret871798"`
	RememberMe bool   `json:"remember_me"`
}
