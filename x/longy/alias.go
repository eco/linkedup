package longy

import (
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = types.ModuleName

	// StoreKey is the key used to access the store
	StoreKey = ModuleName

	// RouterKey is the key for routing messages to our handler
	RouterKey = types.RouterKey
)

var (
	// ModuleCdc is the alias for the amino with the module's types registered
	ModuleCdc = types.ModuleCdc

	// RegisterCodec is the function alias to register types
	RegisterCodec = types.RegisterCodec

	// NewKeeper is the new keeper function alias for longy
	NewKeeper = keeper.NewKeeper

	// NewAttendee is the function alias for creating a new attendee
	NewAttendee = types.NewAttendee

	// NewRekeyMsg is the function alias for the RekeyMsg type
	NewRekeyMsg = types.NewMsgRekey

	// NewQuerier is the fucntion alias for creating a new querier
	NewQuerier = keeper.NewQuerier
)

type (
	// Keeper is the keeper alias for longy
	Keeper = keeper.Keeper

	// Attendee is the type alias for Attendee
	Attendee = types.Attendee

	// RekeyMsg is the type alias for the Rekey Message
	RekeyMsg = types.MsgRekey
)
