package services

import (
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"

	"github.com/chesta132/goreply/reply"
)

func (s *ContextedClass) GetFormTeacher(payload payloads.RequestGetClass) (*models.User, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	var user models.User
	err := s.userRepo.DB().Preload("TeacherProfile").
		Joins("JOIN teachers ON teachers.user_id = users.id").
		Joins("JOIN classes ON classes.form_teacher_id = teachers.id").
		Where("classes.id = ?", payload.ID).
		First(&user).Error

	if err != nil {
		return nil, errorlib.MakeNotFound(err, "class not found", nil)
	}

	return &user, nil
}
