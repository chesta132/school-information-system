package errorlib

import (
	"errors"
	"school-information-system/internal/libs/replylib"

	"github.com/chesta132/goreply/reply"
	"gorm.io/gorm"
)

func MakeUserByTargetIDNotFound(err error) *reply.ErrorPayload {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &reply.ErrorPayload{
			Code:    replylib.CodeNotFound,
			Message: "user not found",
			Fields:  reply.FieldsError{"target_id": "user with targeted id doesn't exist"},
		}
	}
	return &reply.ErrorPayload{
		Code:    replylib.CodeServerError,
		Message: err.Error(),
	}
}

func MakeNotFound(err error, message string, fields reply.FieldsError) *reply.ErrorPayload {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &reply.ErrorPayload{
			Code:    replylib.CodeNotFound,
			Message: message,
			Fields:  fields,
		}
	}
	return &reply.ErrorPayload{
		Code:    replylib.CodeServerError,
		Message: err.Error(),
	}
}

func MakeServerError(err error) *reply.ErrorPayload {
	return &reply.ErrorPayload{
		Code:    replylib.CodeServerError,
		Message: err.Error(),
	}
}
