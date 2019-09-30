package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrefixKey(t *testing.T) {
	prefix := []byte("foo")
	key := []byte("bar")

	prefixedKey := prefixKey(prefix, key)
	require.Equal(t, prefixedKey, []byte("foo::bar"))
}
