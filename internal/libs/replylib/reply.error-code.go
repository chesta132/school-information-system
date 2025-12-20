package replylib

import "net/http"

const (
	CodeNotFound            = "NOT_FOUND"
	CodeServerError         = "SERVER_ERROR"
	CodeBadRequest          = "BAD_REQUEST"
	CodeBadGateway          = "BAD_GATEWAY"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeConflict            = "CONFLICT"
	CodeForbidden           = "FORBIDDEN"
	CodeUnprocessableEntity = "UNPROCESSABLE_ENTITY"
	CodeTooManyRequests     = "TOO_MANY_REQUESTS"
	CodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
	CodeGatewayTimeout      = "GATEWAY_TIMEOUT"
	CodeMethodNotAllowed    = "METHOD_NOT_ALLOWED"
	CodeNotAcceptable       = "NOT_ACCEPTABLE"
	CodeRequestTimeout      = "REQUEST_TIMEOUT"
	CodePayloadTooLarge     = "PAYLOAD_TOO_LARGE"
	CodeUnsupportedMedia    = "UNSUPPORTED_MEDIA_TYPE"
	CodeGone                = "GONE"
	CodeNotImplemented      = "NOT_IMPLEMENTED"
)

var CodeAliases = map[string]int{
	CodeNotFound:            http.StatusNotFound,
	CodeServerError:         http.StatusInternalServerError,
	CodeBadRequest:          http.StatusBadRequest,
	CodeUnauthorized:        http.StatusUnauthorized,
	CodeBadGateway:          http.StatusBadGateway,
	CodeConflict:            http.StatusConflict,
	CodeForbidden:           http.StatusForbidden,
	CodeUnprocessableEntity: http.StatusUnprocessableEntity,
	CodeTooManyRequests:     http.StatusTooManyRequests,
	CodeServiceUnavailable:  http.StatusServiceUnavailable,
	CodeGatewayTimeout:      http.StatusGatewayTimeout,
	CodeMethodNotAllowed:    http.StatusMethodNotAllowed,
	CodeNotAcceptable:       http.StatusNotAcceptable,
	CodeRequestTimeout:      http.StatusRequestTimeout,
	CodePayloadTooLarge:     http.StatusRequestEntityTooLarge,
	CodeUnsupportedMedia:    http.StatusUnsupportedMediaType,
	CodeGone:                http.StatusGone,
	CodeNotImplemented:      http.StatusNotImplemented,
}

func GetCodeByStatus(status int) (code string, ok bool) {
	for k, v := range CodeAliases {
		if v == status {
			return k, true
		}
	}
	return CodeServerError, false
}
