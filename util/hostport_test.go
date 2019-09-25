package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHostAndPort(t *testing.T) {
	goodCases := []string{
		"host:80",
	}
	badCases := []string{
		"",
		"host:80:80",
		":80",
		"host:port",
	}

	for _, server := range goodCases {
		_, _, err := HostAndPort(server)
		require.NoErrorf(t, err, "case: %s", server)
	}
	for _, server := range badCases {
		_, _, err := HostAndPort(server)
		require.Errorf(t, err, "case: %s", server)
	}
}
