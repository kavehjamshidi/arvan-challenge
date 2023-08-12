package handler

const (
	MsgSuccess         = "success"
	MsgValidationError = "validation error"
	MsgFailed          = "failed"
	MsgForbidden       = "forbidden"
	MsgTooManyRequests = "too many requests"
)

type response[T string | *string] struct {
	Error   T      `json:"error"`
	Message string `json:"message"`
}
