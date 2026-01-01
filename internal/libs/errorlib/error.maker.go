package errorlib

import (
	"errors"
	"school-information-system/internal/libs/replylib"
	"strings"

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

func MakeUpdateParentErr(phoneNumber, email string) *reply.ErrorPayload {
	// msg & fields for better error message
	msg := make([]string, 0, 2)
	fields := make([]string, 0, 2)
	if phoneNumber != "" {
		msg = append(msg, "phone number")
		fields = append(fields, "phone")
	}
	if email != "" {
		msg = append(msg, "email")
		fields = append(fields, "email")
	}

	errPayload := &reply.ErrorPayload{
		Code:    replylib.CodeConflict,
		Message: "contact already registered",
		Fields:  reply.FieldsError{},
	}
	for _, field := range fields {
		errPayload.Fields[field] = strings.Join(msg, " or ") + " already registered"
	}

	return errPayload
}
