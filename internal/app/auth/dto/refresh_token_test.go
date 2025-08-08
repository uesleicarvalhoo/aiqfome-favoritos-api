package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto/fixture"
)

func TestRefreshTokenParams_Validate(t *testing.T) {
	t.Parallel()

	builder := fixture.AnyRefreshTokenParams()

	testCases := []struct {
		about         string
		params        dto.RefreshTokenParams
		expectedError string
	}{
		{
			about:         "when refreshToken is empty",
			params:        builder.WithRefreshToken("").Build(),
			expectedError: "[AQF002] refreshToken: campo obrigatório",
		},
		{
			about:         "when refreshToken is provided",
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
