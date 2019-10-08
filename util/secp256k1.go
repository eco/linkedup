package util

import (
	"encoding/hex"
	"fmt"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	secp256k1PrivKeyLen = 32
)

// Secp256k1FromHex parses the hex-encoded `key` string
func Secp256k1FromHex(key string) (tmcrypto.PrivKey, error) {
	bytes, err := hex.DecodeString(TrimHex(key))
	if err != nil {
		return nil, fmt.Errorf("hex decoding: %s", err)
	} else if len(bytes) != secp256k1PrivKeyLen {
		return nil, fmt.Errorf("invalid key byte length. expected: %d, got: %d",
			secp256k1PrivKeyLen, len(bytes))
	}

	var privateKey [secp256k1PrivKeyLen]byte
	copied := copy(privateKey[:], bytes)
	if copied != secp256k1PrivKeyLen {
		errMsg := fmt.Sprintf("expected to copy over %d bytes into the secp256k1 private key",
			secp256k1PrivKeyLen)
		panic(errMsg)
	}

	return secp256k1.PrivKeySecp256k1(privateKey), nil
}
