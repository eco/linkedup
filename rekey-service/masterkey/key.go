package masterkey

import (
	"encoding/hex"
	"fmt"
	"github.com/eco/longy/util"
	"github.com/sirupsen/logrus"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	secp256k1PrivKeyLen = 32
)

var log = logrus.WithField("module", "masterkey")

// Key encapslates the master key for the
// longey game
type Key struct {
	privKey tmcrypto.PrivKey
}

// NewMasterKey is the constructor for `Key`. A new secp256k1 is generated if empty.
// The `chainID` is used when generating RekeyTransactions to prevent cross-chain replay attacks
func NewMasterKey(privateKey tmcrypto.PrivKey) (Key, error) {
	if privateKey == nil {
		key := Key{
			privKey: secp256k1.GenPrivKey(),
		}

		return key, nil
	}

	k := Key{
		privKey: privateKey,
	}

	return k, nil
}

// Secp256k1FromHex parses the hex-encoded `key` string
func Secp256k1FromHex(key string) (tmcrypto.PrivKey, error) {
	if len(key) == 0 {
		log.Info("provided key is empty. generating a new Secp256k1 key")
		return secp256k1.GenPrivKey(), nil
	}

	bytes, err := hex.DecodeString(util.TrimHex(key))
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

// RekeyTransaction generates a `RekeyMsg`, authorized by the master key. The transaction bytes
// generated are created using the cosmos-sdk/x/auth module's StdSignDoc. The account and sequence number
// for the master key is zero.
func (k Key) RekeyTransaction(id int, publicKey []byte) ([]byte, error) {
	return nil, nil
}
