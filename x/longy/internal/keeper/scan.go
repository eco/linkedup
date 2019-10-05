package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

//GetScanByID returns the scan by its id. Returns an error if it cannot find the scan with that id
//nolint:gocritic
func (k Keeper) GetScanByID(ctx sdk.Context, id []byte) (scan types.Scan, err sdk.Error) {
	bz, e := k.Get(ctx, id)
	if e != nil {
		if e.Code() == types.ItemNotFound {
			err = types.ErrScanNotFound("invalid key passed for scan %s", id)
			return
		}
		err = e
		return
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &scan)
	return
}

//SetScan puts the scan into the store with its id as key
//nolint:gocritic
func (k Keeper) SetScan(ctx sdk.Context, scan *types.Scan) {
	k.Set(ctx, scan.ID, k.cdc.MustMarshalBinaryBare(&scan))
}
