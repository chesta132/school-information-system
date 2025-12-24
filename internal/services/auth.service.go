package services

import (
	"context"
	"errors"
	"net/http"
	"school-information-system/config"
	"school-information-system/database/seeds"
	"school-information-system/internal/libs/authlib"
	"school-information-system/internal/libs/phonelib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payload"
	"school-information-system/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Auth struct {
	userRepo    *repos.User
	revokedRepo *repos.Revoked
	adminRepo   *repos.Admin
}

type ContextedAuth struct {
	*Auth
	c   *gin.Context
	ctx context.Context
}

func NewAuth(userRepo *repos.User, revokedRepo *repos.Revoked, adminRepo *repos.Admin) *Auth {
	return &Auth{userRepo, revokedRepo, adminRepo}
}

func (s *Auth) ApplyContext(c *gin.Context) *ContextedAuth {
	return &ContextedAuth{s, c, c.Request.Context()}
}

func (s *ContextedAuth) SignUp(payload payload.RequestSignUp) (*models.User, []http.Cookie, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, nil, errPayload
	}

	// format and validate phone number
	formattedNumber, validNum := phonelib.FormatNumber(payload.Phone)
	if !validNum {
		return nil, nil, &replylib.ErrInvalidPhone
	}

	// check email
	if exists, err := s.userRepo.Exists(s.ctx, "email = ?", payload.Email); err != nil {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	} else if exists {
		return nil, nil, &replylib.ErrEmailRegistered
	}

	// check phone number
	if exists, err := s.userRepo.Exists(s.ctx, "phone = ?", formattedNumber); err != nil {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	} else if exists {
		return nil, nil, &replylib.ErrPhoneRegistered
	}

	// create user object
	hashedPassword, err := authlib.HashPassword(payload.Password)
	if err != nil {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	}

	newUser := &models.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     models.RoleUnsetted,
		Gender:   payload.Gender,
		Phone:    formattedNumber,
	}

	// save user
	if err := s.userRepo.Create(s.ctx, newUser); err != nil {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	}

	// create cookies and return
	access := authlib.CreateAccessCookie(newUser.ID, string(newUser.Role), payload.RememberMe)
	refresh := authlib.CreateRefreshCookie(newUser.ID, string(newUser.Role), payload.RememberMe)

	return newUser, []http.Cookie{access, refresh}, nil
}

func (s *ContextedAuth) SignIn(payload payload.RequestSignIn) (*models.User, []http.Cookie, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, nil, errPayload
	}

	// get user by email
	user, err := s.userRepo.GetFirst(s.ctx, "email = ?", payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, &replylib.ErrEmailNotRegistered
		}
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	}

	// validate password
	if !authlib.ComparePassword(payload.Password, user.Password) {
		return nil, nil, &replylib.ErrIncorrectPassword
	}

	// create cookies and return
	access := authlib.CreateAccessCookie(user.ID, string(user.Role), payload.RememberMe)
	refresh := authlib.CreateRefreshCookie(user.ID, string(user.Role), payload.RememberMe)

	return &user, []http.Cookie{access, refresh}, nil
}

func (s *ContextedAuth) SignOut() []http.Cookie {
	clientRefresh, _ := s.c.Cookie(config.REFRESH_TOKEN_KEY)
	refresh, err := authlib.ParseRefreshToken(clientRefresh)

	// revoke refresh token
	if err == nil {
		revokedToken := &models.Revoked{
			Token:        clientRefresh,
			Reason:       authlib.ReasonUserSignOut,
			RevokedUntil: refresh.ExpiresAt.Time,
		}
		s.revokedRepo.Create(s.ctx, revokedToken)
	}

	acccessCookie := authlib.Invalidate(config.ACCESS_TOKEN_KEY)
	refreshCookie := authlib.Invalidate(config.REFRESH_TOKEN_KEY)
	return []http.Cookie{acccessCookie, refreshCookie}
}

func (s *ContextedAuth) InitiateAdmin(payload payload.RequestInitiateAdmin) (*models.User, *reply.ErrorPayload) {
	// validate payload
	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		return nil, errPayload
	}

	// get and check if targeted user exist
	user, err := s.userRepo.GetByID(s.ctx, payload.TargetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &reply.ErrorPayload{
				Code:    replylib.CodeNotFound,
				Message: "user with targeted id doesn't exist",
				Fields:  []string{"target_id"},
			}
		}
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
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

func (s *ContextedAuth) initiateAdminInTx(payload payload.RequestInitiateAdmin, user models.User) (admin *models.Admin, err error) {
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
