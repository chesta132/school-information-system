package services

import (
	"context"
	"errors"
	"net/http"
	"school-information-system/config"
	"school-information-system/internal/libs/authlib"
	"school-information-system/internal/libs/phonelib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/slicelib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payload"
	"school-information-system/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Auth struct {
	userRepo    *repos.User
	revokedRepo *repos.Revoked
}

type ContextedAuth struct {
	*Auth
	c   *gin.Context
	ctx context.Context
}

func NewAuth(userRepo *repos.User, revokedRepo *repos.Revoked) *Auth {
	return &Auth{userRepo, revokedRepo}
}

func (s *Auth) ApplyContext(c *gin.Context) *ContextedAuth {
	return &ContextedAuth{s, c, c.Request.Context()}
}

func (s *ContextedAuth) SignUp(payload payload.RequestSignUp) (*models.User, []http.Cookie, *reply.ErrorPayload) {
	// validate payload
	if err := validatorlib.Client.Struct(payload); err != nil {
		err, valErrs := validatorlib.TranslateError(err)
		fields := slicelib.Unique(slicelib.Map(valErrs, func(i int, val validator.FieldError) string { return val.Field() }))

		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: "invalid payload",
			Details: err.Error(),
			Fields:  fields,
		}
	}

	// format and validate phone number
	formattedNumber, validNum := phonelib.FormatNumber(payload.Phone)
	if !validNum {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: "invalid phone number",
			Fields:  []string{"phone"},
		}
	}

	// check email
	if exists, err := s.userRepo.Exists(s.ctx, "email = ?", payload.Email); err != nil {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	} else if exists {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeConflict,
			Message: "email already registered",
			Fields:  []string{"email"},
		}
	}

	// check phone number
	if exists, err := s.userRepo.Exists(s.ctx, "phone = ?", formattedNumber); err != nil {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	} else if exists {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeConflict,
			Message: "phone number already registered",
			Fields:  []string{"phone"},
		}
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
	if err := validatorlib.Client.Struct(payload); err != nil {
		err, valErrs := validatorlib.TranslateError(err)
		fields := slicelib.Unique(slicelib.Map(valErrs, func(i int, val validator.FieldError) string { return val.Field() }))

		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: "invalid payload",
			Details: err.Error(),
			Fields:  fields,
		}
	}

	// get user by email
	user, err := s.userRepo.GetFirst(s.ctx, "email = ?", payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, &reply.ErrorPayload{
				Code:    replylib.CodeNotFound,
				Message: "email not registered yet",
				Fields:  []string{"email"},
			}
		}
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	}

	// validate password
	if !authlib.ComparePassword(payload.Password, user.Password) {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeUnauthorized,
			Message: "password is incorrect",
			Fields:  []string{"password"},
		}
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
