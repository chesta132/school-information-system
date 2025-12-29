package services

import (
	"context"
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
	userRepo  *repos.User
	adminRepo *repos.Admin
}

type ContextedAdmin struct {
	*Admin
	c   *gin.Context
	ctx context.Context
}

func NewAdmin(userRepo *repos.User, adminRepo *repos.Admin) *Admin {
	return &Admin{userRepo, adminRepo}
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
	admin, err := s.initiateAdminInTx(payload, &user)
	if err != nil {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	}

	user.AdminProfile = admin
	return &user, nil
}

func (s *ContextedAdmin) initiateAdminInTx(payload payloads.RequestInitiateAdmin, user *models.User) (admin *models.Admin, err error) {
	err = s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		userRepo := s.userRepo.WithTx(tx)
		adminRepo := s.adminRepo.WithTx(tx)

		// update role
		err := userRepo.UpdateByID(s.ctx, user.ID, models.User{Role: models.RoleAdmin})
		if err != nil {
			return err
		}
		user.Role = models.RoleAdmin

		// create admin with permission seeds
		admin = &models.Admin{
			StaffRole:  payload.StaffRole,
			EmployeeID: payload.EmployeeID,
			UserID:     user.ID,
			TimestampJoinTime: models.TimestampJoinTime{
				JoinedAt: payload.JoinedAt,
			},
		}
		err = adminRepo.Create(s.ctx, admin)
		if err != nil {
			return err
		}
		return tx.Model(admin).Association("Permissions").Append(&seeds.PermissionSeeds)
	})
	return
}
