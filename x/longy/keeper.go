package longy

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// AttendeeKeeper maintains the link to data storage and exposes getter/setter
// methods for the various parts of the state machine
type AttendeeKeeper struct {
	coinKeeper   bank.Keeper  //Reference to the bank keeper we use play and pay out rewards
	scanStoreKey sdk.StoreKey // Key for the scan KVStore
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.0
}

// NewKeeper creates new instances of the button Keeper
func NewKeeper(coinKeeper bank.Keeper, scanKey sdk.StoreKey, cdc *codec.Codec) AttendeeKeeper {
	return AttendeeKeeper{
		coinKeeper:   coinKeeper,
		scanStoreKey: scanKey,
		cdc:          cdc,
	}
}
