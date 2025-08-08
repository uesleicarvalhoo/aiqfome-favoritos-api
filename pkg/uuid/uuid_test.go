package uuid_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

func TestNextID(t *testing.T) {
	t.Parallel()

	t.Run("shouldn't be empty", func(t *testing.T) {
		t.Parallel()

		// Action
		sut := uuid.NextID()

		// Assert
		assert.NotEqual(t, uuid.ID{}, sut)
	})

	t.Run("shouldn't equal to last generated ID", func(t *testing.T) {
		t.Parallel()

		// Arrange
		last := uuid.NextID()

		// Arrange
		current := uuid.NextID()

		// Assert
		assert.NotEqual(t, last, current)
	})
}

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		about         string
		stringID      string
		expectedID    uuid.ID
		expectedError error
	}{
		{
			about:         "when ID is invalid",
			stringID:      "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			expectedID:    uuid.ID{},
			expectedError: errors.New("invalid UUID format"),
		},
		{
			about:         "when ID is valid",
			stringID:      "bc8b65bc-08c7-4687-85d3-875b6e0ec449",
			expectedID:    uuid.MustParse("bc8b65bc-08c7-4687-85d3-875b6e0ec449"),
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Action
			parsed, err := uuid.ParseID(tc.stringID)

			// Assert
			assert.Equal(t, tc.expectedID, parsed)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestToString(t *testing.T) {
	t.Parallel()

	// Arrange
	sut := uuid.MustParse("bc8b65bc-08c7-4687-85d3-875b6e0ec449")
	expected := "bc8b65bc-08c7-4687-85d3-875b6e0ec449"

	// Action
	strID := sut.String()

	// Assert
	assert.Equal(t, expected, strID)
}
