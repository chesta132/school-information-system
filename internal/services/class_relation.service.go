package services

import (
	"fmt"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/slicelib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"

	"github.com/chesta132/goreply/reply"
	"gorm.io/gorm"
)

func (s *ContextedClass) GetFormTeacher(payload payloads.RequestGetClass) (*models.User, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	teacher, err := s.classRepo.GetFormTeacher(s.ctx, payload.ID)
	if err != nil {
		return nil, errorlib.MakeNotFound(err, "class not found", nil)
	}

	return &teacher, nil
}

func (s *ContextedClass) GetStudents(payload payloads.RequestGetClass) ([]models.User, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	students, err := s.classRepo.GetStudents(s.ctx, payload.ID)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}

	return students, nil
}

func (s *ContextedClass) GetFull(payload payloads.RequestGetClass) (full *payloads.ResponseGetFullClass, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}
	// prevent nil error
	full = new(payloads.ResponseGetFullClass)

	// get class
	class, err := s.classRepo.GetByID(s.ctx, payload.ID)
	if err != nil {
		return nil, errorlib.MakeNotFound(err, "class not found", nil)
	}
	class.Name = class.GetName()
	full.Class = &class

	// get form teacher
	teacher, err := s.classRepo.GetFormTeacher(s.ctx, payload.ID)
	if err != nil {
		return nil, errorlib.MakeNotFound(err, "form teacher not found", nil)
	}
	full.FormTeacher = &teacher

	// get students
	students, err := s.classRepo.GetStudents(s.ctx, payload.ID)
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}
	full.Students = students

	return
}

func (s *ContextedClass) SetFormTeacher(payload payloads.RequestSetFormTeacher) (teacher *models.User, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	s.classRepo.DB().Transaction(func(tx *gorm.DB) error {
		classRepo := s.classRepo.WithTx(tx)

		// prevent nil error
		teacher = new(models.User)

		// get classes to validate class exists and form teacher is a form teacher of another class
		classes, err := classRepo.GetAll(s.ctx, "id = ? OR form_teacher_id = ?", payload.ID, payload.TeacherID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		var classExist bool
		for _, class := range classes {
			if class.ID == payload.ID {
				classExist = true
			}
			if class.FormTeacherID == payload.TeacherID {
				errPayload = &reply.ErrorPayload{
					Code:    replylib.CodeConflict,
					Message: "teacher is a form teacher in " + class.GetName(),
					Fields: reply.FieldsError{
						"teacher_id": "teacher of this id is a form teacher in " + class.GetName(),
					},
				}
			}
		}
		if !classExist {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "class not found", nil)
			return gorm.ErrRecordNotFound
		}
		if errPayload != nil {
			return gorm.ErrDuplicatedKey
		}

		// get user with teacher profile and validate is teacher exists
		err = tx.Preload("TeacherProfile").
			Joins("JOIN teachers ON teachers.user_id = users.id").
			Where("teachers.id = ?", payload.TeacherID).
			First(teacher).Error
		if err != nil {
			errPayload = errorlib.MakeNotFound(err, "teacher not found", reply.FieldsError{
				"teacher_id": "teacher with this id not found",
			})
			return err
		}

		// update form teacher
		err = classRepo.UpdateByID(s.ctx, payload.ID, models.Class{FormTeacherID: payload.TeacherID})
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}

		return err
	})
	return
}

func (s *ContextedClass) AddStudents(payload payloads.RequestAddStudents) (students []models.User, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	s.classRepo.DB().Transaction(func(tx *gorm.DB) error {
		classRepo := s.classRepo.WithTx(tx)
		studentRepo := s.studentRepo.WithTx(tx)

		// validate class id
		classExists, err := classRepo.Exists(s.ctx, "id = ?", payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !classExists {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "class not found", nil)
			return gorm.ErrRecordNotFound
		}

		// removes duplicate ids
		payload.StudentIDs = slicelib.Unique(payload.StudentIDs)

		// get students and validate length
		err = tx.Preload("StudentProfile").
			Joins("JOIN students ON students.user_id = users.id").
			Where("students.id IN ?", payload.StudentIDs).
			Find(&students).Error
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		notFound := len(payload.StudentIDs) - len(students)
		if notFound > 0 {
			errPayload = &reply.ErrorPayload{
				Code:    replylib.CodeNotFound,
				Message: "student(s) not found",
				Fields: reply.FieldsError{
					"student_ids": fmt.Sprintf("%d student(s) with these id not found", notFound),
				},
			}
			return gorm.ErrRecordNotFound
		}

		// update student's class_id column
		err = studentRepo.Update(s.ctx, models.Student{ClassID: payload.ID}, "id IN ?", payload.StudentIDs)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})
	return
}
