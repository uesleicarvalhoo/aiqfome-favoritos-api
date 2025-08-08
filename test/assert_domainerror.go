package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
)

func AssertDomainError(t *testing.T, err error, expected error) {
	if de, ok := err.(*domainerror.Error); ok {
		assert.Equal(t, expected.(*domainerror.Error).Code, de.Code, "Expected code to be equal")
		assert.Equal(t, expected.(*domainerror.Error).Message, de.Message, "Expected message to be equal")
		for k, v := range de.Details {
			assert.Equal(t, expected.(*domainerror.Error).Details[k], v, "Expected detail to be equal")
		}
	}
}

func AssertValidationDomainError(t *testing.T, err error, expectedErr, expectedField string, errorCode domainerror.ErrorCode) {
	// Unwrap error
	var domainErr *domainerror.Error
	require.True(t, errors.As(err, &domainErr), "expected error to be of type *domainerror.Error")

	// Check error message
	assert.Contains(t, domainErr.Error(), expectedErr)

	// Check specific field and ErrorCode
	if domainErr.Details["fields"] != nil {
		fields, ok := domainErr.Details["fields"].([]validator.ValidationErrorProps)
		require.True(t, ok, "expected fields to be []ValidationErrorProps")

		found := false
		for _, field := range fields {
			if field.Context == expectedField {
				found = true
				assert.Equal(t, errorCode, "unexpected ErrorCode for field %s", field)
			}
		}

		require.Truef(t, found, "expected field '%s' not found in error details. Got fields: %+v", expectedField, fields)
	}
}
