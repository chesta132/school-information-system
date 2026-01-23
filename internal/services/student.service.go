package services

import (
	"context"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/slicelib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Student struct {
	studentRepo *repos.Student
	parentRepo  *repos.Parent
}

type ContextedStudent struct {
	*Student
	c   *gin.Context
	ctx context.Context
}

func NewStudent(studentRepo *repos.Student, parentRepo *repos.Parent) *Student {
	return &Student{studentRepo, parentRepo}
}

func (s *Student) ApplyContext(c *gin.Context) *ContextedStudent {
	return &ContextedStudent{s, c, c.Request.Context()}
}

func (s *ContextedStudent) UpdateStudent(payload payloads.RequestUpdateStudent) (student *models.Student, errPayload *reply.ErrorPayload) {
	payload.ParentIDs = slicelib.Unique(payload.ParentIDs)
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// prevent nil error
	student = new(models.Student)

	s.studentRepo.DB().Transaction(func(tx *gorm.DB) error {
		parentRepo := s.parentRepo.WithTx(tx)
		studentRepo := s.studentRepo.WithTx(tx)

		// check exist
		exists, err := studentRepo.Exists(s.ctx, "id = ?", payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !exists {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "student profile not found", nil)
			return gorm.ErrRecordNotFound
		}

		// update student profile
		if payload.NISN != "" {
			err := studentRepo.Update(s.ctx, models.Student{NISN: payload.NISN}, "id = ?", payload.ID)
			if err != nil {
				errPayload = errorlib.MakeServerError(err)
				return err
			}
		}
		if len(payload.ParentIDs) > 0 {
			parents, err := parentRepo.GetByIDs(s.ctx, payload.ParentIDs)
			if err != nil {
				errPayload = errorlib.MakeServerError(err)
				return err
			}
			if len(parents) != 2 {
				errPayload = &reply.ErrorPayload{Code: replylib.CodeNotFound, Message: "parent(s) not found"}
				return gorm.ErrRecordNotFound
			}
			err = tx.Model(student).Association("Parents").Replace(parents)
			if err != nil {
				errPayload = errorlib.MakeServerError(err)
				return err
			}
		}

		// get to return
		*student, err = studentRepo.GetFirstWithPreload(s.ctx, []string{"Parents"}, "id = ?", payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})

	return
}
