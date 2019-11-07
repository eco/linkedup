package querier

import (
	"encoding/hex"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

//nolint:gocritic
func queryScan(ctx sdk.Context, path []string, keeper keeper.Keeper) (res []byte, err sdk.Error) {
	scan, err := keeper.GetScanByID(ctx, types.Decode(path[0]))
	if err != nil {
		return
	}
	dst := make([]byte, hex.EncodedLen(len(scan.ID)))
	hex.Encode(dst, scan.ID)
	scan.ID = dst
	res, e := codec.MarshalJSONIndent(keeper.Cdc, scan)
	if e != nil {
		panic("could not marshal result to JSON")
	}
	return
}

//nolint:gocritic
func queryScans(ctx sdk.Context, keeper keeper.Keeper) (res []byte, err sdk.Error) {
	scans := keeper.GetAllScans(ctx)
	type ScanString struct {
		ID          string         `json:"id"`
		S1          sdk.AccAddress `json:"s1"`
		S2          sdk.AccAddress `json:"s2"`
		D1          []byte         `json:"d1"`
		D2          []byte         `json:"d2"`
		P1          uint           `json:"p1"`
		P2          uint           `json:"p2"`
		UnixTimeSec int64          `json:"unixTimeSec"`
		Accepted    bool           `json:"accepted"`
	}
	ss := make([]ScanString, len(scans))
	for i := range scans {
		ss[i] = ScanString{
			ID:          types.Encode(scans[i].ID),
			S1:          scans[i].S1,
			S2:          scans[i].S2,
			D1:          scans[i].D1,
			D2:          scans[i].D2,
			P1:          scans[i].P1,
			P2:          scans[i].P2,
			UnixTimeSec: scans[i].UnixTimeSec,
			Accepted:    scans[i].Accepted,
		}
	}
	res, e := json.Marshal(ss)
	//res, e := codec.MarshalJSONIndent(keeper.Cdc, scans)
	if e != nil {
		panic("could not marshal result to JSON")
	}
	err = nil
	return
}
