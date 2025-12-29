package services

import (
	"context"
	"errors"
	"fmt"
	"school-information-system/config"
	"school-information-system/database/seeds"
	"school-information-system/internal/libs/errorlib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/slicelib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payloads"
	"school-information-system/internal/repos"
	"strings"

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
			errPayload = &replylib.ErrPermissionNameExist
			return errors.New(errPayload.Message)
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

func (s *ContextedPermission) GetPermission(payload payloads.RequestGetPermission) (permission *models.Permission, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload = validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return
	}

	perm, err := s.permissionRepo.GetByID(s.ctx, payload.ID)
	if err != nil {
		errPayload = errorlib.MakeNotFound(err, "permission not found", nil)
	}
	permission = &perm
	return
}

func (s *ContextedPermission) GetPermissions(payload payloads.RequestGetPermissions) (permissions []models.Permission, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload = validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return
	}

	q := gorm.G[models.Permission](s.userRepo.DB()).Limit(config.LIMIT_PAGINATED_DATA + 1)
	// action query
	if len(payload.Actions) > 0 {
		payload.Actions = slicelib.Unique(payload.Actions)
		// loop like because actions col is JSON as TEXT
		// its safe cz payload.Actions already validated
		orConditions := make([]string, len(payload.Actions))
		orArgs := make([]any, len(payload.Actions))
		for i, act := range payload.Actions {
			orConditions[i] = "actions LIKE ?"
			orArgs[i] = "%\"" + act + "\"%"
		}
		q = q.Where(strings.Join(orConditions, " OR "), orArgs...)
	}
	// resource query
	if len(payload.Resources) > 0 {
		q = q.Where("resource IN ?", payload.Resources)
	}
	// name match query
	if payload.Query != "" {
		q = q.Where("LOWER(name) LIKE LOWER(?)", "%"+payload.Query+"%")
	}
	// offset query
	if payload.Offset != 0 {
		q = q.Offset(payload.Offset)
	}

	permissions, err := q.Find(s.ctx)
	if err != nil {
		errPayload = errorlib.MakeServerError(err)
	}
	return
}

func (s *ContextedPermission) UpdatePermission(payload payloads.RequestUpdatePermission) (permission *models.Permission, errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload = validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return
	}

	s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		permissionRepo := s.permissionRepo.WithTx(tx)

		// check if permission is permission seeds
		if seeds.IsPermissionSeed(&models.Permission{Id: models.Id{ID: payload.ID}}) {
			errPayload = &replylib.ErrPermissionImmutable
			return errors.New(errPayload.Message)
		}

		// check name uniques and permission presence
		perms, err := gorm.G[models.Permission](tx).Select("id", "name").Where("name = ? OR id = ?", payload.Name, payload.ID).Limit(2).Find(s.ctx)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}

		var permExist bool
		var nameConflict bool

		for _, perm := range perms {
			if perm.ID == payload.ID {
				permExist = true
			}
			// pass if permission.ID = payload.ID
			if payload.Name != "" && perm.Name == payload.Name && perm.ID != payload.ID {
				nameConflict = true
			}
		}

		if !permExist {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeNotFound, Message: "permission not found"}
			return gorm.ErrRecordNotFound
		}

		if nameConflict {
			errPayload = &replylib.ErrPermissionNameExist
			return errors.New(errPayload.Message)
		}

		// update permission
		permission = &models.Permission{Name: payload.Name, Description: payload.Description}
		perm, err := permissionRepo.UpdateByIDAndGet(s.ctx, payload.ID, *permission)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
		} else {
			permission = &perm
		}
		return err
	})
	return
}

func (s *ContextedPermission) DeletePermission(payload payloads.RequestDeletePermission) (errPayload *reply.ErrorPayload) {
	// validate payload
	if errPayload = validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return
	}

	s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		permissionRepo := s.permissionRepo.WithTx(tx)

		// check if permission is permission seeds
		if seeds.IsPermissionSeed(&models.Permission{Id: models.Id{ID: payload.ID}}) {
			errPayload = &replylib.ErrPermissionImmutable
			return errors.New(errPayload.Message)
		}

		// check is permission have relation with another admin
		var m2mExist bool
		err := tx.Raw("SELECT EXISTS (SELECT 1 FROM admin_permissions WHERE permission_id = ? LIMIT 1)", payload.ID).Scan(&m2mExist).Error // line 113
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if m2mExist {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeConflict, Message: "permission still granted by other admin(s)"}
			return errors.New("can't delete granted permission")
		}

		// delete permission
		exists, err := permissionRepo.DeleteByID(s.ctx, payload.ID)
		if err != nil {
			errPayload = errorlib.MakeServerError(err)
			return err
		}
		if !exists {
			errPayload = &reply.ErrorPayload{Code: replylib.CodeNotFound, Message: "permission not found"}
			return gorm.ErrRecordNotFound
		}
		return nil
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
			errPayload = &reply.ErrorPayload{Code: replylib.CodeBadRequest, Message: errorlib.ErrTargetHavePermission.Error()}
			return errorlib.ErrTargetHavePermission
		}

		// get permission to grant
		perm, err := permissionRepo.GetByID(s.ctx, payload.PermissionID)
		if err != nil {
			errPayload = errorlib.MakeNotFound(
				gorm.ErrRecordNotFound,
				"permission not found",
				reply.FieldsError{"permission_id": "permission with this ID not found"},
			)
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
			errPayload = &reply.ErrorPayload{Code: replylib.CodeBadRequest, Message: errorlib.ErrTargetDoesntHavePerm.Error()}
			return errorlib.ErrTargetDoesntHavePerm
		}

		// make sure if permission is seed, it has another granted admin
		if seeds.IsPermissionSeed(permission) {
			selfVal := s.c.MustGet("user")
			self, _ := selfVal.(models.User)

			var exists bool
			err = tx.Raw(
				"SELECT EXISTS (SELECT 1 FROM admin_permissions WHERE admin_id != ? AND permission_id = ? LIMIT 1)",
				self.AdminProfile.ID,
				permission.ID,
			).Scan(&exists).Error
			if err != nil {
				errPayload = errorlib.MakeServerError(err)
				return err
			}

			if !exists {
				err = fmt.Errorf("%w to revoke", errorlib.ErrPermHaventAnotherAdmin)
				errPayload = &reply.ErrorPayload{Code: replylib.CodeConflict, Message: err.Error()}
				return err
			}
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
