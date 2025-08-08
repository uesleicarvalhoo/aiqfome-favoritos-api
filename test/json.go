package test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Marshal(t *testing.T, v any) []byte {
	d, err := json.Marshal(v)
	require.NoError(t, err)

	return d
}
