package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper is a creator for `Keeper`
func NewKeeper(scanKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey: scanKey,
		cdc:      cdc,
	}
}
