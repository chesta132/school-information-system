package services

import (
	"context"
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

type Permission struct {
	userRepo       *repos.User
	permissionRepo *repos.Permission
}

type ContextedPermission struct {
	*Permission
	c   *gin.Context
	ctx context.Context
}

func NewPermission(userRepo *repos.User, permissionRepo *repos.Permission) *Permission {
	return &Permission{userRepo, permissionRepo}
}

func (s *Permission) ApplyContext(c *gin.Context) *ContextedPermission {
	return &ContextedPermission{s, c, c.Request.Context()}
}

func (s *ContextedPermission) GrantPermission(payload payloads.RequestGrantPermission) (user *models.User, permission *models.Permission, errPayload *reply.ErrorPayload) {
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, nil, errPayload
	}

	// transaction to rollback if error
	s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		userRepo := s.userRepo.WithTx(tx)
		permissionRepo := s.permissionRepo.WithTx(tx)

		// get user and validate
		u, err := userRepo.GetFirstWithPreload(s.ctx, []string{"AdminProfile"}, "id = ?", payload.TargetID)
		if err != nil {
			errPayload = errorlib.MakeUserByTargetIDNotFound(err)
			return err
		}
		if u.AdminProfile == nil || u.Role != models.RoleAdmin {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeUnprocessableEntity, Message: errorlib.ErrTargetInvalidRole.Error()}
			return errorlib.ErrTargetInvalidRole
		}
		user = &u

		// get permission to append
		perm, err := permissionRepo.GetByID(s.ctx, payload.PermissionID)
		if err != nil {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "permission not found", []string{"permission_id"})
			return err
		}
		permission = &perm

		// append permission
		err = tx.Model(user.AdminProfile).Association("Permissions").Append(&perm)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})
	return
}
