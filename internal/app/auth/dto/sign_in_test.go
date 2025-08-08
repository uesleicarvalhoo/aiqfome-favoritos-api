package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto/fixture"
)

func TestSignInParams_Validate(t *testing.T) {
	t.Parallel()

	builder := fixture.AnySignInParams()

	testCases := []struct {
		about         string
		params        dto.SignInParams
		expectedError string
	}{
		{
			about:         "when email is invalid",
			params:        builder.WithEmail("not-an-email").Build(),
			expectedError: "[AQF002] email: email inválido",
		},
		{
			about:         "when password is empty",
			params:        builder.WithPassword("").Build(),
			expectedError: "[AQF002] password: campo obrigatório",
		},
		{
			about:         "when both email and password are invalid",
			params:        builder.WithEmail("bad").WithPassword("").Build(),
			expectedError: "[AQF002] email: email inválido; password: campo obrigatório",
		},
		{
			about:         "when all values are valid",
			params:        builder.Build(),
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			err := tc.params.Validate()

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				return
			}
			assert.NoError(t, err)
		})
	}
}
