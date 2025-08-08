package favorite_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/favorite"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about         string
		clientID      uuid.ID
		productID     int
		expectedError string
	}{
		{
			about:         "when clientID is invalid",
			clientID:      uuid.Nil,
			productID:     1,
			expectedError: "[AQF002] clientId: campo obrigatório",
		},
		{
			about:         "when productID is invalid",
			clientID:      uuid.NextID(),
			productID:     0,
			expectedError: "[AQF002] productId: campo obrigatório",
		},
		{
			about:         "when both clientID and productID are invalid",
			clientID:      uuid.Nil,
			productID:     0,
			expectedError: "[AQF002] clientId: campo obrigatório; productId: campo obrigatório",
		},
		{
			about:         "when all values are valid",
			clientID:      uuid.NextID(),
			productID:     42,
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Action
			res, err := favorite.New(tc.clientID, tc.productID)

			// Assert
			if tc.expectedError != "" {
				assert.Equal(t, favorite.Favorite{}, res)
				assert.EqualError(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.clientID, res.ClientID)
			assert.Equal(t, tc.productID, res.ProductID)
			assert.WithinDuration(t, time.Now(), res.RegistredAt, time.Second*1)
		})
	}
}
