package role_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uesleicarvalhoo/aiqfome/role"
)

func TestAction_Validate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about    string
		action   role.Action
		expected bool
	}{
		{
			about:    "when is action read",
			action:   role.ActionRead,
			expected: true,
		},
		{
			about:    "when is action write",
			action:   role.ActionWrite,
			expected: true,
		},
		{
			about:    "when is action delete",
			action:   role.ActionDelete,
			expected: true,
		},
		{
			about:    "when is action manage",
			action:   role.ActionManage,
			expected: true,
		},
		{
			about:    "when action is invalid",
			action:   role.Action(""),
			expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			got := tc.action.IsValid()

			assert.Equal(t, tc.expected, got)
		})
	}
}
