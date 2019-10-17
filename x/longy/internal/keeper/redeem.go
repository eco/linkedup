package keeper

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/tendermint/tendermint/crypto"
)

//SetRedeemAccount sets the redeem account from the genesis file
//nolint:gocritic
func (k Keeper) SetRedeemAccount(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	if addr.Empty() {
		return sdk.ErrInvalidAddress(addr.String())
	}

	account := k.accountKeeper.GetAccount(ctx, addr)
	if account == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	key := types.RedeemKey()
	bz := addr.Bytes()
	k.Set(ctx, key, bz)

	return nil
}

//IsRedeemAccount returns true if the the account passed in is the redeemer account
//nolint:gocritic
func (k Keeper) IsRedeemAccount(ctx sdk.Context, addr sdk.Address) bool {
	key := types.RedeemKey()
	bz, err := k.Get(ctx, key)
	if err != nil {
		return false
	}
	redeemer := sdk.AccAddress(bz)
	return redeemer.Equals(addr)
}

//RedeemPrizes sets all of the prizes for an attendee to claimed = true
//nolint:gocritic
func (k *Keeper) RedeemPrizes(ctx sdk.Context, attendeeAddr sdk.AccAddress) sdk.Error {
	//get the address for the scanned qr code
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

//ValidateSig validates that the hex encoded sig is the signed sha512(badgeId) by the corresponding private key
func ValidateSig(key crypto.PubKey, badgeID string, sig string) sdk.Error {
	if key == nil {
		return types.ErrInvalidSignature("public key cannot be empty")
	}
	hash, e := Hash([]byte(badgeID))
	if e != nil {
		return types.ErrHashingError("error on hashing badge id")
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
