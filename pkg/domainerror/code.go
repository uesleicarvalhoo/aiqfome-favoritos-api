package domainerror

import (
	"net/http"
)

type ErrorCode string

const (
	ErrorCodeNone       ErrorCode = ""
	Default             ErrorCode = "AQF001"
	InvalidParams       ErrorCode = "AQF002"
	ResourceNotFound    ErrorCode = "AQF003"
	DependecyError      ErrorCode = "AQF004"
	OperationNotAllowed ErrorCode = "AQF005"

	// Client
	EmailAlreadyExists ErrorCode = "USR001"
	UserNotActive      ErrorCode = "USR002"

	// Auth
	InvalidPassword       ErrorCode = "AUT001"
	AutenticationNotFound ErrorCode = "AUT002"
	AutenticationInvalid  ErrorCode = "AUT003"

	// Favorites
	ProductAlreadyIsFavorite ErrorCode = "FAV001"
)

func (ec ErrorCode) String() string {
	return string(ec)
}

// NOTE: Remember to keep it sync with ErrorCodes
var httpStatusByCodes = map[ErrorCode]int{
	// General
	ErrorCodeNone:       http.StatusInternalServerError,
	Default:             http.StatusInternalServerError,
	InvalidParams:       http.StatusUnprocessableEntity,
	ResourceNotFound:    http.StatusNotFound,
	DependecyError:      http.StatusInternalServerError,
	OperationNotAllowed: http.StatusForbidden,

	// Client
	EmailAlreadyExists: http.StatusConflict,
	UserNotActive:      http.StatusForbidden,

	// Auth
	InvalidPassword:       http.StatusUnauthorized,
	AutenticationNotFound: http.StatusUnauthorized,
	AutenticationInvalid:  http.StatusUnauthorized,

	// Favorites
	ProductAlreadyIsFavorite: http.StatusConflict,
}

func StatusCode(code ErrorCode) int {
	if sc, ok := httpStatusByCodes[code]; ok {
		return sc
	}

	return http.StatusInternalServerError
}
