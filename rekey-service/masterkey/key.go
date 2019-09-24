package masterkey

import (
	"encoding/hex"
	"fmt"
	"github.com/eco/longy/util"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
)

type privateKey tmcrypto.PrivKeySecp256k1

const (
	privKeyByteLen = 32
)

// Key encapslates the master key for the
// longey game
type Key struct {
	privKey privateKey
}

// NewMasterKey is the constructor for `Key`
func NewMasterKey(hexStr string) (Key, error) {
	hexStr = util.TrimHex(hexStr)
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return Key{}, fmt.Errorf("hex: %s", err)
	} else if len(bytes) != privKeyByteLen {
		return Key{}, fmt.Errorf("incorrect byte length")
	}

	var key privateKey
	copied := copy(key[:], bytes)
	if copied != 32 {
		panic(fmt.Sprintf("key construction %d copy failed", privKeyByteLen))
	}

	k := Key{
		privKey: key,
	}

	return k, nil
}

func (k Key) CreateSignature() string {
	return ""
}
