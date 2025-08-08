package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

func TestAuthClaims_Validate(t *testing.T) {
	t.Parallel()

	claimsBuilder := fixture.AnyAuthClaims()

	testCases := []struct {
		about         string
		claims        dto.AuthClaims
		expectedError string
	}{
		{
			about:         "when userID is zero",
			claims:        claimsBuilder.WithUserID(uuid.Nil).Build(),
			expectedError: "[AQF002] userId: ID do usuário é obrigatório",
		},
		{
			about:         "when userID is valid",
			claims:        claimsBuilder.Build(),
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			err := tc.claims.Validate()

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
