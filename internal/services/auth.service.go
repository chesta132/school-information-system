package services

import (
	"context"
	"errors"
	"net/http"
	"school-information-system/internal/libs/authlib"
	"school-information-system/internal/libs/replylib"
	"school-information-system/internal/libs/slicelib"
	"school-information-system/internal/libs/validatorlib"
	"school-information-system/internal/models"
	"school-information-system/internal/models/payload"
	"school-information-system/internal/repos"
	"strings"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Auth struct {
	userRepo *repos.User
}

type ContextedAuth struct {
	*Auth
	c   *gin.Context
	ctx context.Context
}

func NewAuth(userRepo *repos.User) *Auth {
	return &Auth{userRepo}
}

func (s *Auth) ApplyContext(c *gin.Context) *ContextedAuth {
	return &ContextedAuth{s, c, c.Request.Context()}
}

func (s *ContextedAuth) SignIn(payload payload.RequestSignIn) (*models.User, []http.Cookie, *reply.ErrorPayload) {
	if err := validatorlib.Client.Struct(payload); err != nil {
		err, valErrs := validatorlib.TranslateError(err)
		fields := slicelib.Map(valErrs, func(i int, val validator.FieldError) string { return val.Field() })
		fields = slicelib.Unique(fields)

		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: "invalid payload",
			Details: err.Error(),
			Field:   strings.Join(fields, ", "),
		}
	}

	user, err := s.userRepo.GetFirst(s.ctx, "email = ?", payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, &reply.ErrorPayload{
				Code:    replylib.CodeNotFound,
				Message: "email not registered yet",
				Field:   "email",
			}
		}
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeServerError,
			Message: err.Error(),
		}
	}

	if !authlib.ComparePassword(payload.Password, user.Password) {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeUnauthorized,
			Message: "password is incorrect",
			Field:   "password",
		}
	}

	access := authlib.CreateAccessCookie(user.ID, string(user.Role), payload.RememberMe)
	refresh := authlib.CreateRefreshCookie(user.ID, string(user.Role), payload.RememberMe)

	return &user, []http.Cookie{access, refresh}, nil
}
