package internal

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to data storage and exposes getter/setter
// methods for the various parts of the state machine
type Keeper struct {
	accountKeeper *auth.AccountKeeper //The keeper for managing accounts
	bankKeeper    *bank.Keeper        // Reference to the bank keeper that we will use to prevent spamming if we have time
	attendeeStore AttendeeStore       // Store for the scan KVStore
	cdc           *codec.Codec        // The wire codec for binary encoding/decoding.0
}

// NewKeeper creates new instances of the button Keeper
func NewKeeper(accountKeeper *auth.AccountKeeper, bankKeeper *bank.Keeper, attendeeStoreKey sdk.StoreKey,
	cdc *codec.Codec) Keeper {
	attendeeStore := NewAttendeeStore(attendeeStoreKey, cdc)
	return Keeper{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		attendeeStore: attendeeStore,
		cdc:           cdc,
	}
}

//GetAccountKeeper returns the account keeper
func (k *Keeper) GetAccountKeeper() *auth.AccountKeeper {
	return k.accountKeeper
}

//GetBankKeeper returns the bank keeper
func (k *Keeper) GetBankKeeper() *bank.Keeper {
	return k.bankKeeper
}

//GetAttendeeStore returns the store for managing attendee state
func (k *Keeper) GetAttendeeStore() AttendeeStore {
	return k.attendeeStore
}
