package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/user"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
)

func TestFromDomain(t *testing.T) {
	t.Parallel()

	id := uuid.NextID()
	userBuilder := fixtureUser.AnyUser().
		WithID(id).
		WithName("Ueslei Carvalho").
		WithEmail("ueslei@email.com").
		WithActive(true)

	clientBuilder := fixture.AnyClient().
		WithID(id).
		WithName("Ueslei Carvalho").
		WithEmail("ueslei@email.com").
		WithActive(true)

	tests := []struct {
		about    string
		user     user.User
		expected dto.Client
	}{
		{
			about:    "when everything is fine",
			user:     userBuilder.Build(),
			expected: clientBuilder.Build(),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			// Arrange

			// Action
			got := dto.FromDomain(tc.user)

			// Assert
			assert.Equal(t, tc.expected, got)
		})
	}
}
