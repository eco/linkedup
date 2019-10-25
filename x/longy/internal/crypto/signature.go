package crypto

import (
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

	_, e := sdk.AccAddressFromBech32(address)
	if e != nil {
		return sdk.ErrInvalidAddress("could not decode bech32 address")
	}

	sigBytes, e := hex.DecodeString(sig)
	if e != nil {
		return types.ErrSigDecodeError("error on hex decoding signature")
	}

	if !key.VerifyBytes([]byte(address), sigBytes) {
		return types.ErrInvalidSignature("signature does not match account public key")
	}

	return nil
}
