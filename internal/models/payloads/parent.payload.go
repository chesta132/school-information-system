package payloads

import "school-information-system/internal/models"

type RequestCreateParent struct {
	FullName string            `json:"full_name" validate:"required" example:"Chesta Ardiona"`
	Phone    string            `json:"phone" validate:"required" example:"+6281234567890"`
	Email    string            `json:"email" validate:"required,email" example:"chestaardi4@gmail.com"`
	Gender   models.UserGender `json:"gender" validate:"required,user_gender"`
}

type RequestGetParent struct {
	ID string `uri:"id" validate:"required,uuid4"`
}

type RequestGetParents struct {
	Offset int               `form:"offset"`
	Query  string            `form:"q"`
	Gender models.UserGender `form:"gender"`
	Email  string            `form:"email"`
}

type RequestUpdateParent struct {
	ID       string            `uri:"id" validate:"required,uuid4"`
	FullName string            `json:"full_name" example:"Chesta Ardiona"`
	Phone    string            `json:"phone" example:"+6281234567890"`
	Email    string            `json:"email" validate:"omitempty,email" example:"chestaardi4@gmail.com"`
	Gender   models.UserGender `json:"gender" validate:"omitempty,user_gender"`
}
