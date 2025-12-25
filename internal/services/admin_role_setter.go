package services

import (
	"context"
	"errors"
	"fmt"
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

type RoleSetter struct {
	userRepo    *repos.User
	adminRepo   *repos.Admin
	studentRepo *repos.Student
	classRepo   *repos.Class
	parentRepo  *repos.Parent
	teacherRepo *repos.Teacher
	subjectRepo *repos.Subject
}

type ContextedRoleSetter struct {
	*RoleSetter
	c   *gin.Context
	ctx context.Context
}

func NewRoleSetter(userRepo *repos.User, adminRepo *repos.Admin, studentRepo *repos.Student, classRepo *repos.Class, parentRepo *repos.Parent, teacherRepo *repos.Teacher, subjectRepo *repos.Subject) *RoleSetter {
	return &RoleSetter{userRepo, adminRepo, studentRepo, classRepo, parentRepo, teacherRepo, subjectRepo}
}

func (s *RoleSetter) ApplyContext(c *gin.Context) *ContextedRoleSetter {
	return &ContextedRoleSetter{s, c, c.Request.Context()}
}

func (s *ContextedRoleSetter) SetRole(payload payloads.RequestSetRole) (*models.User, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// get and check user
	user, err := s.userRepo.GetByID(s.ctx, payload.TargetID)
	if err != nil {
		return nil, errorlib.MakeUserByTargetIDNotFound(err)
	}

	switch payload.TargetRole {
	case models.RoleStudent:
		// set role student
		var payload payloads.RequestSetRoleStudent
		if err := s.c.ShouldBindJSON(&payload); err != nil {
			return nil, &reply.ErrorPayload{Code: replylib.CodeBadRequest, Message: err.Error()}
		}
		stud, err := s.setRoleStudent(payload, user)
		user.StudentProfile = stud
		return &user, err

		// set role teacher
	case models.RoleTeacher:
		var payload payloads.RequestSetRoleTeacher
		if err := s.c.ShouldBindJSON(&payload); err != nil {
			return nil, &reply.ErrorPayload{Code: replylib.CodeBadRequest, Message: err.Error()}
		}
		teacher, err := s.setRoleTeacher(payload, user)
		user.TeacherProfile = teacher
		return &user, err

		// set role admin
	case models.RoleAdmin:
		var payload payloads.RequestSetRoleAdmin
		if err := s.c.ShouldBindJSON(&payload); err != nil {
			return nil, &reply.ErrorPayload{Code: replylib.CodeBadRequest, Message: err.Error()}
		}
		admin, err := s.setRoleAdmin(payload, user)
		user.AdminProfile = admin
		return &user, err
	}

	return nil, &reply.ErrorPayload{Code: replylib.CodeBadRequest, Message: "invalid target role", Fields: []string{"target_role"}}
}

func (s *ContextedRoleSetter) setRoleStudent(payload payloads.RequestSetRoleStudent, user models.User) (student *models.Student, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// transaction to rollback if error
	s.adminRepo.DB().Transaction(func(tx *gorm.DB) error {
		classRepo := s.classRepo.WithTx(tx)
		parentRepo := s.parentRepo.WithTx(tx)
		studentRepo := s.studentRepo.WithTx(tx)

		// check is class exists
		classExists, err := classRepo.Exists(s.ctx, "id = ?", payload.ClassID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !classExists {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "class not found", []string{"class_id"})
			return gorm.ErrRecordNotFound
		}

		// check is found parent 2
		parents, err := parentRepo.GetByIDs(s.ctx, payload.ParentIDs)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if len(parents) != 2 {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeConflict, Message: "existing parents must be 2"}
			return errors.New("found parent is not 2")
		}

		// create student
		student = &models.Student{
			NISN:    payload.NISN,
			UserID:  user.ID,
			ClassID: payload.ClassID,
		}
		err = studentRepo.Create(s.ctx, student)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}

		// append parents to student
		err = tx.Model(student).Association("Parents").Append(parents)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})
	return
}

func (s *ContextedRoleSetter) setRoleTeacher(payload payloads.RequestSetRoleTeacher, user models.User) (teacher *models.Teacher, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// transaction to rollback if error
	s.adminRepo.DB().Transaction(func(tx *gorm.DB) error {
		subjectRepo := s.subjectRepo.WithTx(tx)
		teacherRepo := s.teacherRepo.WithTx(tx)

		// check if all subjects exists
		subjects, err := subjectRepo.GetByIDs(s.ctx, payload.SubjectIDs)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if subjectIDsLen, subjectsLen := len(payload.SubjectIDs), len(subjects); subjectIDsLen > subjectsLen {
			notFound := subjectIDsLen - subjectsLen
			err := fmt.Errorf("%d subject(s) not found", notFound)
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, err.Error(), []string{"subject_ids"})
			return err
		}

		// check is unique conflict
		teacherExist, err := teacherRepo.Exists(s.ctx, "nuptk = ? OR employee_id = ?", payload.NUPTK, payload.EmployeeID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if teacherExist {
			err := errors.New("other teacher with same NUTPK or employee id already exist")
			errPayload = &reply.ErrorPayload{
				Code:    replylib.CodeConflict,
				Message: err.Error(),
				Fields:  []string{"nuptk", "employee_id"},
			}
			return err
		}

		// create teacher
		teacher = &models.Teacher{
			NUPTK:      payload.NUPTK,
			EmployeeID: payload.EmployeeID,
			JoinedAt:   payload.JoinedAt,
			UserID:     user.ID,
		}
		err = teacherRepo.Create(s.ctx, teacher)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}

		// append subjects to teacher
		err = tx.Model(teacher).Association("Subjects").Append(subjects)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})
	return
}

func (s *ContextedRoleSetter) setRoleAdmin(payload payloads.RequestSetRoleAdmin, user models.User) (admin *models.Admin, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// transaction to rollback if error
	s.adminRepo.DB().Transaction(func(tx *gorm.DB) error {
		adminRepo := s.adminRepo.WithTx(tx)

		// check is unique conflict
		adminExist, err := adminRepo.Exists(s.ctx, "employee_id = ?", payload.EmployeeID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if adminExist {
			err := errors.New("other admin with same employee id already exist")
			errPayload = &reply.ErrorPayload{
				Code:    replylib.CodeConflict,
				Message: err.Error(),
				Fields:  []string{"employee_id"},
			}
			return err
		}

		// create admin
		// set permissions with other route
		admin = &models.Admin{
			StaffRole:  payload.StaffRole,
			EmployeeID: payload.EmployeeID,
			JoinedAt:   payload.JoinedAt,
			UserID:     user.ID,
		}
		return adminRepo.Create(s.ctx, admin)
	})
	return
}
