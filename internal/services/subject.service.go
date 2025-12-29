package services

import (
	"context"
	"log"
	"school-information-system/config"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (s *ContextedSubject) GetSubject(payload payloads.RequestGetSubject) (*models.Subject, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// get subject
	subject, err := s.subjetRepo.GetByID(s.ctx, payload.ID)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	return &subject, nil
}

func (s *ContextedSubject) GetSubjects(payload payloads.RequestGetSubjects) ([]models.Subject, *reply.ErrorPayload) {
	log.Println(payload)
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// get subjects
	q := gorm.G[models.Subject](s.subjetRepo.DB()).
		Limit(config.LIMIT_PAGINATED_DATA + 1)
	if payload.Offset > 0 {
		q = q.Offset(payload.Offset)
	}
	if payload.Query != "" {
		q = q.Where("LOWER(name) LIKE LOWER(?)", "%"+payload.Query+"%")
	}

	subject, err := q.Find(s.ctx)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	return subject, nil
}
