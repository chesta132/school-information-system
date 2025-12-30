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

type Class struct {
	classRepo   *repos.Class
	userRepo    *repos.User
	teacherRepo *repos.Teacher
	studentRepo *repos.Student
}

type ContextedClass struct {
	*Class
	c   *gin.Context
	ctx context.Context
}

func NewClass(classRepo *repos.Class, userRepo *repos.User, teacherRepo *repos.Teacher, studentRepo *repos.Student) *Class {
	return &Class{classRepo, userRepo, teacherRepo, studentRepo}
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

	class, err := s.classRepo.GetByID(s.ctx, payload.ID)
	if err != nil {
		return nil, errorlib.MakeNotFound(err, "class not found", nil)
	}

	class.GetName()
	return &class, nil
}

func (s *ContextedClass) GetClasses(payload payloads.RequestGetClasses) ([]models.Class, *reply.ErrorPayload) {
	q := gorm.G[models.Class](s.classRepo.DB()).Limit(config.LIMIT_PAGINATED_DATA + 1).Offset(payload.Offset)
	if payload.Grade > 0 && payload.Grade <= 12 {
		q = q.Where("grade = ?", payload.Grade)
	}
	if payload.ClassNumber > 0 {
		q = q.Where("class_number = ?", payload.ClassNumber)
	}
	if payload.Major != "" {
		q = q.Where("major = ?", payload.Major)
	}

	classes, err := q.Find(s.ctx)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	for i, class := range classes {
		classes[i].Name = class.GetName()
	}

	return classes, nil
}

func (s *ContextedClass) UpdateClass(payload payloads.RequestUpdateClass) (class *models.Class, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	s.classRepo.DB().Transaction(func(tx *gorm.DB) error {
		classRepo := s.classRepo.WithTx(tx)

		// check presence
		classExists, err := classRepo.Exists(s.ctx, "id = ?", payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !classExists {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "class not found", nil)
			return gorm.ErrRecordNotFound
		}

		// update class
		update := models.Class{Grade: payload.Grade, Major: payload.Major, ClassNumber: payload.ClassNumber}
		classUpdt, err := classRepo.UpdateByIDAndGet(s.ctx, payload.ID, update)

		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		classUpdt.Name = classUpdt.GetName()
		class = &classUpdt
		return nil
	})
	return
}

func (s *ContextedClass) DeleteClass(payload payloads.RequestDeleteClass) (errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return errPayload
	}

	s.classRepo.DB().Transaction(func(tx *gorm.DB) error {
		classRepo := s.classRepo.WithTx(tx)
		studentRepo := s.studentRepo.WithTx(tx)

		// check student relation
		studentExists, err := studentRepo.Exists(s.ctx, "class_id = ?", payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if studentExists {
			err = errors.New("class still related with student(s)")
			errPayload = &reply.ErrorPayload{
				Code:    replylib.CodeConflict,
				Message: err.Error(),
			}
			return err
		}

		// delete class
		ok, err := classRepo.DeleteByID(s.ctx, payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !ok {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "class not found", nil)
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	return
}
