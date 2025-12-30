package services

import (
	"context"
	"errors"
	"school-information-system/config"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Subject struct {
	subjectRepo *repos.Subject
	teacherRepo *repos.Teacher
}

type ContextedSubject struct {
	*Subject
	c   *gin.Context
	ctx context.Context
}

func NewSubject(subjectRepo *repos.Subject, teacherRepo *repos.Teacher) *Subject {
	return &Subject{subjectRepo, teacherRepo}
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
	err := s.subjectRepo.Create(s.ctx, subject)
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
	subject, err := s.subjectRepo.GetByID(s.ctx, payload.ID)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	return &subject, nil
}

func (s *ContextedSubject) GetSubjects(payload payloads.RequestGetSubjects) ([]models.Subject, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// get subjects
	q := gorm.G[models.Subject](s.subjectRepo.DB()).
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

func (s *ContextedSubject) UpdateSubject(payload payloads.RequestUpdateSubject) (*models.Subject, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// update and get subject
	subject := models.Subject{Name: payload.Name}
	subject, err := s.subjectRepo.UpdateByIDAndGet(s.ctx, payload.ID, subject)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	return &subject, nil
}

func (s *ContextedSubject) DeleteSubject(payload payloads.RequestDeleteSubject) (errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return errPayload
	}

	s.subjectRepo.DB().Transaction(func(tx *gorm.DB) error {
		subjectRepo := s.subjectRepo.WithTx(tx)

		// validate if there is another teacher related to this subject
		var related bool
		err := tx.Raw("SELECT EXISTS (SELECT 1 FROM teacher_subjects WHERE subject_id = ? LIMIT 1)", payload.ID).Scan(&related).Error
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if related {
			errPayload = &reply.ErrorPayload{
				Code:    replylib.CodeConflict,
				Message: "subject still registered by other teacher(s)",
			}
			return errors.New("can't delete related subject")
		}

		// delete subject
		ok, err := subjectRepo.DeleteByID(s.ctx, payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !ok {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "subject not found", nil)
			return gorm.ErrRecordNotFound
		}

		return nil
	})

	return
}

func (s *ContextedSubject) GetTeacherOfSubject(payload payloads.RequestGetTeacherOfSubject) (teachers []models.Teacher, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// validate presencce
	exists, err := s.subjectRepo.Exists(s.ctx, "id = ?", payload.SubjectID)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}
	if !exists {
		return nil, errorlib.MakeNotFound(gorm.ErrRecordNotFound, "subject not found", nil)
	}

	// get teachers
	err = s.teacherRepo.DB().WithContext(s.ctx).
		Joins("JOIN teacher_subjects ON teacher_subjects.teacher_id = teachers.id").
		Where("teacher_subjects.subject_id = ?", payload.SubjectID).
		Find(&teachers).Error

	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	return
}
