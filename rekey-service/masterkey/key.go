package masterkey

import (
	"encoding/hex"
	"fmt"
	"github.com/eco/longy/util"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	privateKeyByteLen = 32
)

// Key encapslates the master key for the
// longey game
type Key struct {
	privKey tmcrypto.PrivKeySecp256k1
}

// NewMasterKey is the constructor for `Key`. A new key will be generated hexStr is empty
func NewMasterKey(hexStr string) (Key, error) {
	hexStr = util.TrimHex(hexStr)
	if len(hexStr) == 0 {
		key := Key{
			privKey: tmcrypto.GenPrivKey(),
		}

		return key, nil
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return Key{}, fmt.Errorf("hex: %s", err)
	} else if len(bytes) != privateKeyByteLen {
		return Key{}, fmt.Errorf("incorrect byte length. got %d, expected: %d",
			len(bytes), privateKeyByteLen)
	}

	var key tmcrypto.PrivKeySecp256k1
	copied := copy(key[:], bytes)
	if copied != 32 {
		panic(fmt.Sprintf("key construction %d copy failed", privateKeyByteLen))
	}

	k := Key{
		privKey: key,
	}

	return k, nil
}

// RekeySignature generates the signature signed by the master key allowing
// attendee `id` to reset with the given `nonce`.
//
// The signature is over
// SHA256("resetkey(id=<id>, nonce=<nonce>)")
func (k Key) RekeySignature(id, nonce int) ([]byte, error) {
	bytesToSign := []byte(fmt.Sprintf("resetkey(id=%d, nonce=%d)", id, nonce))
	hash := tmhash.Sum(bytesToSign)

	sig, err := k.privKey.Sign(hash)
	if err != nil {
		return nil, fmt.Errorf("tmcrypto: %s", err)
	}

	return sig, nil
}
