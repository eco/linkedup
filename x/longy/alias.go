package longy

import (
	"github.com/eco/longy/x/longy/internal"
	"github.com/eco/longy/x/longy/internal/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = types.ModuleName
	// RouterKey is the key for routing messages to our handler
	RouterKey = types.RouterKey
	// StoreKey is the key for the keeper store
	StoreKey = types.StoreKey
)

var (
	// NewKeeper is the new keeper function alias for longy
	NewKeeper = internal.NewKeeper

	// NewRekeyMsg is the function alias for the RekeyMsg type
	NewRekeyMsg = types.NewMsgRekey
)

type (
	// Keeper is the keeper alias for longy
	Keeper = internal.Keeper

	// RekeyMsg is the type alias for the Rekey Message
	RekeyMsg = types.MsgRekey
)
