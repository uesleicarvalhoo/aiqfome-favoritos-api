package domainerror

import (
	"fmt"
)

type Error struct {
	Code    ErrorCode
	Message string
	Details map[string]any
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s | cause: %v", e.Code, e.Message, e.Cause)
	}

	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(code ErrorCode, msg string, details map[string]any) error {
	return &Error{
		Code:    code,
		Message: msg,
		Details: details,
	}
}

func Wrap(err error, code ErrorCode, msg string, details map[string]any) error {
	return &Error{
		Code:    code,
		Message: msg,
		Details: details,
		Cause:   err,
	}
}
