package replylib

import (
	"github.com/chesta132/goreply/reply"
)

var (
	// phone number errors

	ErrInvalidPhone = reply.ErrorPayload{
		Code:    CodeBadRequest,
		Message: "invalid phone number",
		Fields:  []string{"phone"},
	}
	ErrPhoneRegistered = reply.ErrorPayload{
		Code:    CodeConflict,
		Message: "phone number already registered",
		Fields:  []string{"phone"},
	}

	// email errors

	ErrEmailRegistered = reply.ErrorPayload{
		Code:    CodeConflict,
		Message: "email already registered",
		Fields:  []string{"email"},
	}
	ErrEmailNotRegistered = reply.ErrorPayload{
		Code:    CodeNotFound,
		Message: "email not registered yet",
		Fields:  []string{"email"},
	}

	// password errors

	ErrIncorrectPassword = reply.ErrorPayload{
		Code:    CodeUnauthorized,
		Message: "password is incorrect",
		Fields:  []string{"password"},
	}
	ErrIncorrectKey = reply.ErrorPayload{
		Code:    CodeUnauthorized,
		Message: "invalid key",
		Fields:  []string{"key"},
	}

	// presence/absence error

	ErrAdminExist = reply.ErrorPayload{
		Code:    CodeConflict,
		Message: "admin is already exist",
	}
	ErrPermissionNameExist = reply.ErrorPayload{
		Code:    CodeConflict,
		Message: "another permission with this name already registered",
		Fields:  []string{"name"},
	}
	ErrPermissionImmutable = reply.ErrorPayload{
		Code:    CodeUnprocessableEntity,
		Message: "this permission is immutable",
	}
)

func ErrorPayloadToArgs(errPayload *reply.ErrorPayload) (code, message string, opt reply.OptErrorPayload) {
	return errPayload.Code, errPayload.Message, reply.OptErrorPayload{Details: errPayload.Details, Fields: errPayload.Fields}
}
