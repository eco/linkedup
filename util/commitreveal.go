package util

import (
	"bytes"
	"github.com/google/uuid"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

// Commitment is the Pedersen Commitmentment
type Commitment []byte

// NewCommitment will create a Commitmentment using `secret`
func NewCommitment(secret []byte) Commitment {
	hash := tmcrypto.Sha256([]byte(secret))

	return Commitment(hash)
}

// CreateCommitment will return a newly generated secret with it's corresponding Commitment
func CreateCommitment() (secret []byte, commitment Commitment) {
	rand := uuid.New()

	secret = rand[:]
	commitment = NewCommitment(secret)

	return secret, commitment
}

// Equals checks if c == c2
func (c Commitment) Equals(c2 Commitment) bool {
	return bytes.Equal(c[:], c2[:])
}

// Empty is an indicator for a nil Commitment
func (c Commitment) Empty() bool {
	return len(c) == 0
}

// Bytes returns the underlying bytes of this commitment
func (c Commitment) Bytes() []byte {
	return c[:]
}

// Len returns the byte length of this commitment {
func (c Commitment) Len() int {
	return len(c.Bytes())
}

// VerifyReveal will verify the Commitmentment against `reveal`
func (c Commitment) VerifyReveal(secret []byte) bool {
	expected := NewCommitment(secret)

	return c.Equals(expected)
}
