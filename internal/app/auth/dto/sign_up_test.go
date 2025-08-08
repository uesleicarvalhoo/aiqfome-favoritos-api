package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto/fixture"
)

func TestSignUpParams_Validate(t *testing.T) {
	t.Parallel()

	builder := fixture.AnySignUpParams()

	testCases := []struct {
		about         string
		params        dto.SignUpParams
		expectedError string
	}{
		{
			about:         "when email is invalid",
			params:        builder.WithEmail("not-an-email").Build(),
			expectedError: "[AQF002] email: email inválido",
		},
		{
			about:         "when name is empty",
			params:        builder.WithName("").Build(),
			expectedError: "[AQF002] nome: campo obrigatório",
		},
		{
			about:         "when password is empty",
			params:        builder.WithPassword("").Build(),
			expectedError: "[AQF002] password: campo obrigatório",
		},
		{
			about:         "when email and name are invalid",
			params:        builder.WithEmail("bad").WithName("").Build(),
			expectedError: "[AQF002] email: email inválido; nome: campo obrigatório",
		},
		{
			about:         "when email and password are invalid",
			params:        builder.WithEmail("bad").WithPassword("").Build(),
			expectedError: "[AQF002] email: email inválido; password: campo obrigatório",
		},
		{
			about:         "when name and password are invalid",
			params:        builder.WithName("").WithPassword("").Build(),
			expectedError: "[AQF002] nome: campo obrigatório; password: campo obrigatório",
		},
		{
			about:         "when all fields are invalid",
			params:        builder.WithEmail("").WithName("").WithPassword("").Build(),
			expectedError: "[AQF002] email: email inválido; nome: campo obrigatório; password: campo obrigatório",
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
