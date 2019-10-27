package longy

import (
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/querier"
	"github.com/eco/longy/x/longy/internal/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = types.ModuleName
	// StoreKey is the key used to access the store
	StoreKey = types.StoreKey
	// RouterKey is the key for routing messages to our handler
	RouterKey = types.RouterKey

	/** ErrCodes **/

	// CodeAttendeeKeyed is the alias for AttendeeKeyed
	CodeAttendeeKeyed = types.AttendeeKeyed
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

	// NewMsgKey is the function alias for the MsgKey type
	NewMsgKey = types.NewMsgKey

	// NewMsgBonus is the function alias for the MsgBonus type
	NewMsgBonus = types.NewMsgBonus

	// NewMsgClearBonus is the function alias for the MsgBonus type
	NewMsgClearBonus = types.NewMsgClearBonus

	// NewQuerier is the function alias for creating a new querier
	NewQuerier = querier.NewQuerier
)

type (
	// Keeper is the keeper alias for longy
	Keeper = keeper.Keeper

	// Attendee is the type alias for Attendee
	Attendee = types.Attendee

	// MsgKey is the type alias for MsgKey
	MsgKey = types.MsgKey

	// GenesisAttendees is the array of attendees for the genesis file
	GenesisAttendees = types.GenesisAttendees

	// GenesisScans is the array of scans for the genesis file
	GenesisScans = types.GenesisScans

	// GenesisPrizes is the array of prizes for the event
	GenesisPrizes = types.GenesisPrizes

	// GenesisService is the genesis type for the service account
	GenesisService = types.GenesisService
)
