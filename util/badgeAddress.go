package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
)

// IDToAddress will deterministically calculate an `sdk.AccAddress` using the given `id`
func IDToAddress(id string) sdk.AccAddress {
	privKey := tmcrypto.GenPrivKeySecp256k1([]byte(id))
	address := privKey.PubKey().Address()

	return sdk.AccAddress(address)
}
