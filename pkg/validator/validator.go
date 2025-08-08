package validator

import (
	"fmt"
	"strings"

	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
)

type ValidationErrorProps struct {
	Context string `json:"context"`
	Message string `json:"message"`
}

type Validator struct {
	errors []ValidationErrorProps
}

func New() Validator {
	return Validator{
		errors: make([]ValidationErrorProps, 0),
	}
}

func (v *Validator) AddError(context string, message string) {
	v.errors = append(v.errors, ValidationErrorProps{
		Context: context,
		Message: message,
	})
}

func (v *Validator) AddSubErrors(prefix string, subErr error) {
	if subErr == nil {
		return
	}

	if domainErr, ok := subErr.(*domainerror.Error); ok {
		fields, ok := domainErr.Details["fields"].([]ValidationErrorProps)
		if ok {
			for _, fieldErr := range fields {
				v.errors = append(v.errors, ValidationErrorProps{
					Context: fmt.Sprintf("%s.%s", prefix, fieldErr.Context),
					Message: fieldErr.Message,
				})
			}
		}
	}
}

func (v Validator) Validate() error {
	if len(v.errors) == 0 {
		return nil
	}

	var msgBuilder strings.Builder
	for _, err := range v.errors {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s; ", err.Context, err.Message))
	}

	return domainerror.New(
		domainerror.InvalidParams,
		strings.TrimSuffix(msgBuilder.String(), "; "),
		map[string]any{
			"fields": v.errors,
		},
	)
}
