package longy

import (
	"github.com/eco/longy/x/longy/internal"
	"github.com/eco/longy/x/longy/internal/types"
)

const (
	//ModuleName is the name of the module
	ModuleName = types.ModuleName
	//RouterKey is the key for routing messages to our handler
	RouterKey = types.RouterKey
	//StoreKey is the key for the keeper store
	StoreKey = types.StoreKey
)

var (
	//NewKeeper is the new keeper function alias for longy
	NewKeeper = internal.NewKeeper
	//RegisterCodec is the codec register
	RegisterCodec = types.RegisterCodec
	//ModuleCdc is the cdc for all the message types
	ModuleCdc = types.ModuleCdc
)

type (
	//Keeper is the keeper alias for longy
	Keeper = internal.Keeper
)
