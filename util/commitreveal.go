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
	hash := tmcrypto.Sha256(secret)

	return Commitment(hash)
}

// CreateCommitment will return a newly generated secret with it's corresponding Commitment
func CreateCommitment() (secret string, commitment Commitment) {
	secret = uuid.New().String()
	commitment = NewCommitment([]byte(secret))

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
func (c Commitment) VerifyReveal(secret string) bool {
	return c.VerifyRevealBytes([]byte(secret))
}

// VerifyRevealBytes will verifty against `reveal`
func (c Commitment) VerifyRevealBytes(secret []byte) bool {
	expected := NewCommitment(secret)
	return c.Equals(expected)
}
