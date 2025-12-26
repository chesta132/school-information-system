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

func (s *ContextedPermission) validateAdminAndPermission(ctx context.Context, tx *gorm.DB, targetID, permissionID string) (user *models.User, hasPermission bool, errPayload *reply.ErrorPayload, err error) {
	userRepo := s.userRepo.WithTx(tx)

	u, err := userRepo.GetFirstWithPreload(ctx, []string{"AdminProfile", "AdminProfile.Permissions"}, "id = ?", targetID)
	if err != nil {
		errPayload = errorlib.MakeUserByTargetIDNotFound(err)
		return
	}
	if u.AdminProfile == nil || u.Role != models.RoleAdmin {
		errPayload = &reply.ErrorPayload{Code: replylib.CodeUnprocessableEntity, Message: errorlib.ErrTargetInvalidRole.Error()}
		err = errorlib.ErrTargetInvalidRole
		return
	}

	user = &u
	for _, perm := range u.AdminProfile.Permissions {
		if perm.ID == permissionID {
			hasPermission = true
			break
		}
	}
	return
}

func (s *ContextedPermission) GrantPermission(payload payloads.RequestGrantPermission) (user *models.User, permission *models.Permission, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, nil, errPayload
	}

	// transaction to rollback if error
	s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		permissionRepo := s.permissionRepo.WithTx(tx)

		// get user and validate
		var hasPermission bool
		var err error
		user, hasPermission, errPayload, err = s.validateAdminAndPermission(s.ctx, tx, payload.TargetID, payload.PermissionID)
		if err != nil {
			return err
		}
		if hasPermission {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeConflict, Message: errorlib.ErrTargetHavePermission.Error()}
			return errorlib.ErrTargetHavePermission
		}

		// get permission to grant
		perm, err := permissionRepo.GetByID(s.ctx, payload.PermissionID)
		if err != nil {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "permission not found", []string{"permission_id"})
			return err
		}
		permission = &perm

		// grant permission
		err = tx.Model(user.AdminProfile).Association("Permissions").Append(&perm)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})
	return
}

func (s *ContextedPermission) RevokePermission(payload payloads.RequestRevokePermission) (user *models.User, permission *models.Permission, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, nil, errPayload
	}

	// transaction to rollback if error
	s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		permissionRepo := s.permissionRepo.WithTx(tx)

		// get user and validate
		var hasPermission bool
		var err error
		user, hasPermission, errPayload, err = s.validateAdminAndPermission(s.ctx, tx, payload.TargetID, payload.PermissionID)
		if err != nil {
			return err
		}
		if !hasPermission {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeUnprocessableEntity, Message: errorlib.ErrTargetDoesntHavePerm.Error()}
			return errorlib.ErrTargetDoesntHavePerm
		}

		// get permission to revoke
		perm, err := permissionRepo.GetByID(s.ctx, payload.PermissionID)
		if err != nil {
			errPayload = errorlib.MakeNotFound(gorm.ErrRecordNotFound, "permission not found", []string{"permission_id"})
			return err
		}
		permission = &perm

		// revoke permission
		err = tx.Model(user.AdminProfile).Association("Permissions").Delete(&perm)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})
	return
}
