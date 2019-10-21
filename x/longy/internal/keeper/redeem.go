package keeper

import (
	"crypto/sha512"
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/tendermint/tendermint/crypto"
)

//RedeemPrizes sets all of the prizes for an attendee to claimed = true
//nolint:gocritic
func (k *Keeper) RedeemPrizes(ctx sdk.Context, attendeeAddr sdk.AccAddress) sdk.Error {
	//get the AccAddress for the scanned qr code
	attendee, ok := k.GetAttendee(ctx, attendeeAddr)
	if !ok {
		return types.ErrAttendeeNotFound("cannot find the attendee")
	}

	winnings := attendee.Winnings
	for i := range winnings {
		winnings[i].Claimed = true
	}

	k.SetAttendee(ctx, &attendee)
	return nil
}

//ValidateSig validates that the hex encoded sig is the signed sha512(address) by the corresponding private key
func ValidateSig(key crypto.PubKey, address string, sig string) sdk.Error {
	if key == nil {
		return types.ErrInvalidSignature("public key cannot be empty")
	}
	hash, e := Hash([]byte(address))
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
