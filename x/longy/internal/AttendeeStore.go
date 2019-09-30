package internal

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/errors"
	"github.com/eco/longy/x/longy/internal/types"
)

const (
	//AttendeeStoreKey is the key for the KVStore
	AttendeeStoreKey = "attendee"
)

// AttendeeStore is the structure for the AttendeeStore class.
type AttendeeStore struct {
	Store
}

// NewAttendeeStore creates a new instance of an AttendeeStore
func NewAttendeeStore(attendeeStoreKey sdk.StoreKey, cdc *codec.Codec) AttendeeStore {
	return AttendeeStore{
		NewDefaultStore(attendeeStoreKey, cdc),
	}
}

// GetAttendee returns the attendee for a given key, or an error if it could not be found. The key is
// the attendee.GetIDBytes()
// nolint: gocritic
func (a *AttendeeStore) GetAttendee(ctx sdk.Context, key []byte) (attendee *types.Attendee, err sdk.Error) {
	item, err := a.GetItemBytes(ctx, key)

	if err != nil {
		if err.Code() == errors.ItemNotFound {
			err = errors.ErrAttendeeNotFound("could not find attendee with id : %s", key)
		}

		return
	}

	a.cdc.MustUnmarshalBinaryBare(item, &attendee)
	return
}

// SetAttendee sets the attendee
// nolint: gocritic
func (a *AttendeeStore) SetAttendee(ctx sdk.Context, attendee *types.Attendee) {
	a.PutItem(ctx, attendee.GetIDBytes(), attendee)
}
