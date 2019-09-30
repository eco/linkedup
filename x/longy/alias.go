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
	// RegisterCodec is the function alias to register types
	RegisterCodec = types.RegisterCodec

	// NewKeeper is the new keeper function alias for longy
	NewKeeper = keeper.NewKeeper

	// NewRekeyMsg is the function alias for the RekeyMsg type
	NewRekeyMsg = types.NewMsgRekey
)

type (
	// Keeper is the keeper alias for longy
	Keeper = keeper.Keeper

	// RekeyMsg is the type alias for the Rekey Message
	RekeyMsg = types.MsgRekey
)
