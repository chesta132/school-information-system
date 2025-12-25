package services

import (
	"context"
	"errors"
	"school-information-system/config"
	"school-information-system/database/seeds"
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

type Admin struct {
	userRepo    *repos.User
	adminRepo   *repos.Admin
	studentRepo *repos.Student
	classRepo   *repos.Class
	parentRepo  *repos.Parent
}

type ContextedAdmin struct {
	*Admin
	c   *gin.Context
	ctx context.Context
}

func NewAdmin(userRepo *repos.User, adminRepo *repos.Admin, studentRepo *repos.Student, classRepo *repos.Class, parentRepo *repos.Parent) *Admin {
	return &Admin{userRepo, adminRepo, studentRepo, classRepo, parentRepo}
}

func (s *Admin) ApplyContext(c *gin.Context) *ContextedAdmin {
	return &ContextedAdmin{s, c, c.Request.Context()}
}

func (s *ContextedAdmin) InitiateAdmin(payload payloads.RequestInitiateAdmin) (*models.User, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// get and check if targeted user exist
	user, err := s.userRepo.GetByID(s.ctx, payload.TargetID)
	if err != nil {
		return nil, errorlib.MakeUserByTargetIDNotFound(err)
	}

	// check is any other admin exist
	if exists, err := s.adminRepo.Exists(s.ctx, "1 = 1"); err != nil {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	} else if exists {
		return nil, &replylib.ErrAdminExist
	}

	// validate key
	if payload.Key != config.INITIATE_ADMIN_KEY {
		return nil, &replylib.ErrIncorrectKey
	}

	// transaction to rollback if error
	admin, err := s.initiateAdminInTx(payload, user)
	if err != nil {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	}

	user.AdminProfile = admin
	return &user, nil
}

func (s *ContextedAdmin) initiateAdminInTx(payload payloads.RequestInitiateAdmin, user models.User) (admin *models.Admin, err error) {
	err = s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		userRepo := s.userRepo.WithTx(tx)
		adminRepo := s.adminRepo.WithTx(tx)

		// update role
		err := userRepo.UpdateByID(s.ctx, user.ID, models.User{Role: models.RoleAdmin})
		if err != nil {
			return err
		}

		// create admin with permission seeds
		admin = &models.Admin{
			StaffRole:  payload.StaffRole,
			EmployeeID: payload.EmployeeID,
			UserID:     user.ID,
			JoinedAt:   payload.JoinedAt,
		}
		err = adminRepo.Create(s.ctx, admin)
		if err != nil {
			return err
		}
		return tx.Model(admin).Association("Permissions").Append(&seeds.PermissionSeeds)
	})
	return
}

func (s *ContextedAdmin) SetRole(payload payloads.RequestSetRole) (*models.User, *reply.ErrorPayload) {
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
		var payload payloads.RequestSetRoleStudent
		if err := s.c.ShouldBindJSON(&payload); err != nil {
			return nil, &reply.ErrorPayload{Code: replylib.CodeBadRequest, Message: err.Error()}
		}
		stud, err := s.setRoleStudent(payload, user)
		user.StudentProfile = stud
		return &user, err
	}

	return nil, nil
}

func (s *ContextedAdmin) setRoleStudent(payload payloads.RequestSetRoleStudent, user models.User) (student *models.Student, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	s.adminRepo.DB().Transaction(func(tx *gorm.DB) error {
		classRepo := s.classRepo.WithTx(tx)
		parentRepo := s.parentRepo.WithTx(tx)
		studentRepo := s.studentRepo.WithTx(tx)

		classExists, err := classRepo.Exists(s.ctx, "id = ?", payload.ClassID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !classExists {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "class not found", []string{"class_id"})
			return gorm.ErrRecordNotFound
		}

		parents, err := parentRepo.GetByIDs(s.ctx, payload.ParentIDs)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if len(parents) != 2 {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeConflict, Message: "existing parents must be 2"}
			return errors.New("found parent is not 2")
		}

		student = &models.Student{
			NISN:    payload.NISN,
			UserID:  payload.TargetID,
			ClassID: payload.ClassID,
		}

		err = studentRepo.Create(s.ctx, student)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}

		err = tx.Model(student).Association("Parents").Append(parents)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})
	return
}
