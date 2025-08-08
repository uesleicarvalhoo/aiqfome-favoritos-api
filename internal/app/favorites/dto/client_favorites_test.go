package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

func TestGetClientFavoritesParams_Validate(t *testing.T) {
	t.Parallel()

	builder := fixture.AnyGetClientFavoritesParams()

	testCases := []struct {
		about         string
		params        dto.GetClientFavoritesParams
		expectedError string
	}{
		{
			about:         "when clientID is zero",
			params:        builder.WithClientID(uuid.Nil).Build(),
			expectedError: "[AQF002] clientId: campo obrigatório",
		},
		{
			about:         "when pageSize is less than 1",
			params:        builder.WithPageSize(0).Build(),
			expectedError: "[AQF002] pageSize: deve ser maior do que 1",
		},
		{
			about:         "when page is negative",
			params:        builder.WithPage(-1).Build(),
			expectedError: "[AQF002] page: não pode ser negativo",
		},
		{
			about:         "when multiple fields are invalid",
			params:        builder.WithClientID(uuid.Nil).WithPageSize(0).WithPage(-1).Build(),
			expectedError: "[AQF002] clientId: campo obrigatório; pageSize: deve ser maior do que 1; page: não pode ser negativo",
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
