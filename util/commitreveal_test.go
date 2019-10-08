package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommitReveal(t *testing.T) {
	s, c := CreateCommitment()
	require.True(t, c.VerifyReveal(s))

	secret := []byte("foo")
	commitment := NewCommitment(secret)
	require.True(t, commitment.VerifyRevealBytes(secret))
}
