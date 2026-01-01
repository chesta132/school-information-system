package services

import (
	"context"
	"fmt"
	"school-information-system/config"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/phonelib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/repos"
	"strings"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (s *ContextedParent) GetParent(payload payloads.RequestGetParent) (*models.Parent, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	parent, err := s.parentRepo.GetByID(s.ctx, payload.ID)
	if err != nil {
		return nil, errorlib.MakeNotFound(err, "parent not found", nil)
	}

	return &parent, nil
}

func (s *ContextedParent) GetParents(payload payloads.RequestGetParents) ([]models.Parent, *reply.ErrorPayload) {
	q := gorm.G[models.Parent](s.parentRepo.DB()).Limit(config.LIMIT_PAGINATED_DATA + 1).Offset(payload.Offset)
	if payload.Email != "" {
		q = q.Where("email = ?", payload.Email)
	}
	if payload.Gender != "" && (payload.Gender == models.GenderFemale || payload.Gender == models.GenderMale) {
		q = q.Where("gender = ?", payload.Gender)
	}
	if payload.Query != "" {
		q = q.Where("full_name LIKE ?", "%"+payload.Query+"%")
	}

	parents, err := q.Find(s.ctx)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}
	return parents, nil
}

func (s *ContextedParent) UpdateParent(payload payloads.RequestUpdateParent) (*models.Parent, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// format and validate phone number
	var formattedNumber string
	if payload.Phone != "" {
		var validNum bool
		formattedNumber, validNum = phonelib.FormatNumber(payload.Phone)
		if !validNum {
			return nil, &replylib.ErrInvalidPhone
		}
	}

	if formattedNumber != "" || payload.Email != "" {
		// build query for validate
		query := make([]string, 0, 2)
		args := []any{payload.ID}
		if formattedNumber != "" {
			query = append(query, "phone = ?")
			args = append(args, formattedNumber)
		}
		if payload.Email != "" {
			query = append(query, "email = ?")
			args = append(args, payload.Email)
		}

		// validate phone number and email
		built := fmt.Sprintf("id != ? AND (%s)", strings.Join(query, " OR "))
		exists, err := s.parentRepo.Exists(s.ctx, built, args...)
		if err != nil {
			return nil, errorlib.MakeServerError(err)
		}
		if exists {
			return nil, errorlib.MakeUpdateParentErr(formattedNumber, payload.Email)
		}
	}

	// update parent, it's ok for let all value enters because gorm auto ignore zero value
	parent := models.Parent{
		FullName: payload.FullName,
		Phone:    formattedNumber,
		Email:    payload.Email,
		Gender:   payload.Gender,
	}

	parent, err := s.parentRepo.UpdateByIDAndGet(s.ctx, payload.ID, parent)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	return &parent, nil
}
