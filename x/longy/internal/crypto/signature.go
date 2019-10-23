package crypto

import (
	"crypto/sha512"
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/tendermint/tendermint/crypto"
)

//ValidateSig validates that the hex encoded sig is the signed sha512(address) by the corresponding private key
func ValidateSig(key crypto.PubKey, address string, sig string) sdk.Error {
	if key == nil {
		return types.ErrInvalidPublicKey("public key cannot be empty")
	}

	accAddr, e := sdk.AccAddressFromBech32(address)
	if e != nil {
		return sdk.ErrInvalidAddress("could not decode bech32 address")
	}
	//hexAddr := hex.EncodeToString(accAddr)
	//fmt.Printf("hex addr : %v", hexAddr)

	hash, e := Hash(accAddr)
	//hexHash := hex.EncodeToString(hash)
	//fmt.Printf("hex hash : %v", hexHash)
	if e != nil {
		return types.ErrHashingError("error on hashing address")
	}

	sigBytes, e := hex.DecodeString(sig)
	if e != nil {
		return types.ErrSigDecodeError("error on hex decoding signature")
	}

	if !key.VerifyBytes(hash, sigBytes) {
		return types.ErrInvalidSignature("signature does not match account public key")
	}

	return nil
}

//Hash hashes a varargs of byte arrays with sha512 in order of first to last
func Hash(toHash ...[]byte) (hash []byte, err error) {
	hasher := sha512.New()

	for i := range toHash {
		_, err = hasher.Write(toHash[i])
		if err != nil {
			return
		}
	}

	return hasher.Sum(nil), err
}
