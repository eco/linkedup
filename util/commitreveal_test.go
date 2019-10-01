package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommitReveal(t *testing.T) {
	secret, commitment := CreateCommitment()
	require.True(t, commitment.VerifyReveal(secret))

	secret = []byte("foo")
	commitment = NewCommitment(secret)
	require.True(t, commitment.VerifyReveal(secret))
}
