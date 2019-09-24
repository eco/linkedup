package masterkey

import (
	"crypto/ecdsa"
)

// Key encapslates the master key for the
// longey game
type Key struct {
	privKey *ecdsa.PrivateKey
}

// NewMaskerKey is the constructor for `Key`
func NewMasterKey(hexStr string) (Key, error) {
	k := Key{
		privKey: nil,
	}

	return k, nil
}
