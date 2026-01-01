package services

import (
	"context"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/phonelib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Parent struct {
	parentRepo *repos.Parent
}

type ContextedParent struct {
	*Parent
	c   *gin.Context
	ctx context.Context
}

func NewParent(parentRepo *repos.Parent) *Parent {
	return &Parent{parentRepo}
}

func (s *Parent) ApplyContext(c *gin.Context) *ContextedParent {
	return &ContextedParent{s, c, c.Request.Context()}
}

func (s *ContextedParent) CreateParent(payload payloads.RequestCreateParent) (*models.Parent, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// format and validate phone number
	formattedNumber, validNum := phonelib.FormatNumber(payload.Phone)
	if !validNum {
		return nil, &replylib.ErrInvalidPhone
	}

	// validate phone number and email
	exists, err := s.parentRepo.Exists(s.ctx, "phone = ? OR email = ?", formattedNumber, payload.Email)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}
	if exists {
		msg := "phone number or email already registered"
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeConflict,
			Message: "contact already registered",
			Fields: reply.FieldsError{
				"email": msg,
				"phone": msg,
			},
		}
	}

	// create parent
	parent := &models.Parent{
		FullName: payload.FullName,
		Phone:    formattedNumber,
		Email:    payload.Email,
		Gender:   payload.Gender,
	}

	err = s.parentRepo.Create(s.ctx, parent)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	return parent, nil
}
