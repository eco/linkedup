package internal

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/errors"
)

// Store is the structure for the Store class.
type Store struct {
	storeKey sdk.StoreKey // Unexposed key to access verifier store from sdk.Context
	cdc      *codec.Codec
}

// StoreItem is
type StoreItem interface {
	GetIDBytes() []byte
}

// NewDefaultStore creates a new instance of an Store
func NewDefaultStore(storeKey sdk.StoreKey, cdc *codec.Codec) Store {
	return Store{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// AddItem adds an item to the store.
func (s Store) AddItem(ctx *sdk.Context, item StoreItem) {
	s.PutItem(ctx, item.GetIDBytes(), item)
}

// PutItem adds an item to the store using the given key.
func (s Store) PutItem(ctx *sdk.Context, key []byte, item interface{}) {
	store := s.getStore(ctx)
	store.Set(key, s.cdc.MustMarshalBinaryBare(item))
}

// GetItemBytes returns an item from the store.
func (s Store) GetItemBytes(ctx *sdk.Context, key []byte) (v []byte, err sdk.Error) {
	store := s.getStore(ctx)
	v = store.Get(key)
	if len(v) == 0 {
		err = errors.ErrItemNotFound("invalid key passed for item %s", key)
		return
	}

	return
}

// ExistsItem returns true if a given dentity exists in the store.
func (s Store) ExistsItem(ctx *sdk.Context, key []byte) bool {
	_, err := s.GetItemBytes(ctx, key)
	return err == nil
}

// RemoveItem removes an item from the store.
func (s Store) RemoveItem(ctx *sdk.Context, key []byte) {
	store := s.getStore(ctx)
	store.Delete(key)
}

func (s Store) getStore(ctx *sdk.Context) sdk.KVStore {
	return ctx.KVStore(s.storeKey)
}

//nolint: unused
func (s Store) getCodec() *codec.Codec {
	return s.cdc
}
