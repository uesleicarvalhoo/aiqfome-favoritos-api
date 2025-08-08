package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

func TestRemoveProductFromFavoritesParams_Validate(t *testing.T) {
	t.Parallel()

	builder := fixture.AnyRemoveProductFromFavoritesParams()

	testCases := []struct {
		about         string
		params        dto.RemoveProductFromFavoritesParams
		expectedError string
	}{
		{
			about:         "when clientID is zero",
			params:        builder.WithClientID(uuid.Nil).Build(),
			expectedError: "[AQF002] clientId: campo obrigatório",
		},
		{
			about:         "when productID is zero",
			params:        builder.WithProductID(0).Build(),
			expectedError: "[AQF002] productId: campo obrigatório",
		},
		{
			about:         "when both clientID and productID are invalid",
			params:        builder.WithClientID(uuid.Nil).WithProductID(0).Build(),
			expectedError: "[AQF002] clientId: campo obrigatório; productId: campo obrigatório",
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
