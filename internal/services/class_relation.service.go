package services

import (
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"

	"github.com/chesta132/goreply/reply"
)

func (s *ContextedClass) getFormTeacher(id string) (models.User, error) {
	var user models.User
	err := s.userRepo.DB().Preload("TeacherProfile").
		Joins("JOIN teachers ON teachers.user_id = users.id").
		Joins("JOIN classes ON classes.form_teacher_id = teachers.id").
		Where("classes.id = ?", id).
		First(&user).Error
	return user, err
}

func (s *ContextedClass) GetFormTeacher(payload payloads.RequestGetClass) (*models.User, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	user, err := s.getFormTeacher(payload.ID)
	if err != nil {
		return nil, errorlib.MakeNotFound(err, "class not found", nil)
	}

	return &user, nil
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
	teacher, err := s.getFormTeacher(payload.ID)
	if err != nil {
		return nil, errorlib.MakeNotFound(err, "form teacher not found", nil)
	}
	full.FormTeacher = &teacher

	// get students
	var students []models.User
	err = s.userRepo.DB().Preload("StudentProfile").
		Joins("JOIN students ON students.user_id = users.id").
		Where("students.class_id = ?", payload.ID).
		Find(&students).Error
	if err != nil {
		return nil, errorlib.MakeServerError(err)
	}
	full.Students = students

	return
}
