package services

import (
	"context"
	"errors"
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

type Class struct {
	classRepo   *repos.Class
	teacherRepo *repos.Teacher
}

type ContextedClass struct {
	*Class
	c   *gin.Context
	ctx context.Context
}

func NewClass(classRepo *repos.Class, teacherRepo *repos.Teacher) *Class {
	return &Class{classRepo, teacherRepo}
}

func (s *Class) ApplyContext(c *gin.Context) *ContextedClass {
	return &ContextedClass{s, c, c.Request.Context()}
}

func (s *ContextedClass) CreateClass(payload payloads.RequestCreateClass) (class *models.Class, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	s.classRepo.DB().Transaction(func(tx *gorm.DB) error {
		classRepo := s.classRepo.WithTx(tx)
		teacherRepo := s.teacherRepo.WithTx(tx)

		// check is teacher exists
		teacherExists, err := teacherRepo.Exists(s.ctx, "id = ?", payload.FormTeacherID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !teacherExists {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "teacher not found", reply.FieldsError{
				"form_teacher_id": "teacher with this id not found",
			})
			return gorm.ErrRecordNotFound
		}

		// check is name exists
		classExists, err := classRepo.Exists(s.ctx,
			"grade = ? AND major = ? AND class_number = ?",
			payload.Grade, payload.Major, payload.ClassNumber,
		)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if classExists {
			msg := "class with this name already exist"
			errPayload = &reply.ErrorPayload{
				Code:    replylib.CodeConflict,
				Message: msg,
				Fields: reply.FieldsError{
					"grade":        msg,
					"major":        msg,
					"class_number": msg,
				},
			}
			return errors.New(msg)
		}

		// create class
		class = &models.Class{
			Grade:         payload.Grade,
			Major:         payload.Major,
			ClassNumber:   payload.ClassNumber,
			FormTeacherID: payload.FormTeacherID,
		}
		class.GetName()

		err = classRepo.Create(s.ctx, class)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}

		return nil
	})

	return
}

func (s *ContextedClass) GetClass(payload payloads.RequestGetClass) (*models.Class, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	class, err := gorm.G[models.Class](s.classRepo.DB()).Preload("FormTeacher", func(db gorm.PreloadBuilder) error {
		db.Select("id")
		return nil
	}).Where("id = ?", payload.ID).First(s.ctx)
	if err != nil {
		return nil, errorlib.MakeNotFound(err, "class not found", nil)
	}

	class.GetName()
	return &class, nil
}
