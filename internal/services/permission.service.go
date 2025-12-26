package services

import (
	"context"
	"errors"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/slicelib"
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

func (s *ContextedPermission) validateAdminAndPermission(ctx context.Context, tx *gorm.DB, targetID, permissionID string) (user *models.User, permission *models.Permission, errPayload *reply.ErrorPayload, err error) {
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
			permission = perm
			break
		}
	}
	return
}

func (s *ContextedPermission) CreatePermission(payload payloads.RequestCreatePermission) (permission *models.Permission, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		permissionRepo := s.permissionRepo.WithTx(tx)

		// validate unique
		if exists, err := permissionRepo.Exists(s.ctx, "name = ?", payload.Name); err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		} else if exists {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeConflict, Message: "permission with this name already exist"}
			return errors.New("name permission is not unique")
		}

		// prevent double
		payload.Actions = slicelib.Unique(payload.Actions)

		// create permission
		permission = &models.Permission{
			Name:        payload.Name,
			Resource:    payload.Resource,
			Description: payload.Description,
			Actions:     payload.Actions,
		}
		err := permissionRepo.Create(s.ctx, permission)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})
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
		var err error
		user, permission, errPayload, err = s.validateAdminAndPermission(s.ctx, tx, payload.TargetID, payload.PermissionID)
		if err != nil {
			return err
		}
		if permission != nil {
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
		// get user & permission then validate
		var err error
		user, permission, errPayload, err = s.validateAdminAndPermission(s.ctx, tx, payload.TargetID, payload.PermissionID)
		if err != nil {
			return err
		}
		if permission == nil {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeUnprocessableEntity, Message: errorlib.ErrTargetDoesntHavePerm.Error()}
			return errorlib.ErrTargetDoesntHavePerm
		}

		// revoke permission
		err = tx.Model(user.AdminProfile).Association("Permissions").Delete(permission)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		}
		return err
	})
	return
}
