package services

import (
	"context"
	"fmt"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Teacher struct {
	teacherRepo *repos.Teacher
	subjectRepo *repos.Subject
}

type ContextedTeacher struct {
	*Teacher
	c   *gin.Context
	ctx context.Context
}

func NewTeacher(teacherRepo *repos.Teacher, subjectRepo *repos.Subject) *Teacher {
	return &Teacher{teacherRepo, subjectRepo}
}

func (s *Teacher) ApplyContext(c *gin.Context) *ContextedTeacher {
	return &ContextedTeacher{s, c, c.Request.Context()}
}

func (s *ContextedTeacher) UpdateTeacher(payload payloads.RequestUpdateTeacher) (teacher *models.Teacher, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// prevent nil error
	teacher = new(models.Teacher)

	s.teacherRepo.DB().Transaction(func(tx *gorm.DB) error {
		teacherRepo := s.teacherRepo.WithTx(tx)
		subjectRepo := s.subjectRepo.WithTx(tx)

		// check exist
		exists, err := teacherRepo.Exists(s.ctx, "id = ?", payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !exists {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "teacher profile not found", nil)
			return gorm.ErrRecordNotFound
		}

		// update profile
		if payload.NUPTK != "" {
			err := teacherRepo.Update(s.ctx, models.Teacher{NUPTK: payload.NUPTK}, "id = ?", payload.ID)
			if err != nil {
				errPayload = errorlib.MakeServerError(err)
				return err
			}
		}
		if len(payload.SubjectIDs) > 0 {
			subjects, err := subjectRepo.GetByIDs(s.ctx, payload.SubjectIDs)
			if err != nil {
				errPayload = errorlib.MakeServerError(err)
				return err
			}
			notFound := len(subjects) - len(payload.SubjectIDs)
			if notFound > 0 {
				errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "subject(s) not found", reply.FieldsError{
					"subject_ids": fmt.Sprintf("%d subject(s) not found", notFound),
				})
				return gorm.ErrRecordNotFound
			}
			err = tx.Model(teacher).Where("id = ?", payload.ID).Association("Subjects").Replace(subjects)
			if err != nil {
				errPayload = errorlib.MakeServerError(err)
				return err
			}
		}

		*teacher, err = teacherRepo.GetByID(s.ctx, payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})

	return
}
