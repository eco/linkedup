package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTrimHex(t *testing.T) {
	expected := "0123"

	cases := []string{
		"0x123",
		"0X123",
		"0x0123",
		"0X0123",
		"0123",
		"123",
	}

	for _, str := range cases {
		got := TrimHex(str)
		require.Equal(t, got, expected)
	}

	require.Equal(t, TrimHex(""), "")
}
