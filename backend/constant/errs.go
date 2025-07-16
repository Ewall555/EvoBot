package constant

import "errors"

const (
	CodeSuccess           = 200
	CodeErrBadRequest     = 400
	CodeErrUnauthorized   = 401
	CodeErrNotFound       = 404
	CodeErrInternalServer = 500
)

var (
	ErrAuth            = errors.New("ErrAuth")
	ErrInitialPassword = errors.New("ErrInitialPassword")
	ErrRecordNotFound  = errors.New("ErrRecordNotFound")
	ErrNotSupportType  = errors.New("ErrNotSupportType")
	ErrInvalidParams   = errors.New("ErrInvalidParams")

	ErrTokenParse = errors.New("ErrTokenParse")
)

var (
	ErrTypeInternalServer = "ErrInternalServer"
)
