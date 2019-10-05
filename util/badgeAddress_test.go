package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIDToAddress(t *testing.T) {
	// sanity check that this is deterministic
	id := "longy"
	address := IDToAddress(id)
	address2 := IDToAddress(id)

	require.Equal(t, address, address2)
}
