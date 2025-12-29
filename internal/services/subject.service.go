package services

import (
	"context"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Subject struct {
	subjetRepo *repos.Subject
}

type ContextedSubject struct {
	*Subject
	c   *gin.Context
	ctx context.Context
}

func NewSubject(subjetRepo *repos.Subject) *Subject {
	return &Subject{subjetRepo}
}

func (s *Subject) ApplyContext(c *gin.Context) *ContextedSubject {
	return &ContextedSubject{s, c, c.Request.Context()}
}

func (s *ContextedSubject) CreateSubject(payload payloads.RequestCreateSubject) (*models.Subject, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// create subject
	subject := &models.Subject{Name: payload.Name}
	err := s.subjetRepo.Create(s.ctx, subject)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	return subject, nil
}
