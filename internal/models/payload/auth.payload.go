package payload

import "school-information-system/internal/models"

type RequestSignUp struct {
	FullName   string            `json:"full_name" validate:"required"`
	Email      string            `json:"email" validate:"required,email"`
	Password   string            `json:"password" validate:"required,min=8"`
	Gender     models.UserGender `json:"gender" validate:"required,oneof=male female"`
	Phone      string            `json:"phone" validate:"required"`
	RememberMe bool              `json:"remember_me"`
}

type RequestSignIn struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}
