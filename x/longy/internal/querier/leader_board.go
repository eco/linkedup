package querier

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
	"sort"
	"time"
)

//LeaderBoard returns the leader board after building it from the attendees in the event
//nolint:gocritic,unparam
func leaderBoard(ctx sdk.Context, keeper keeper.Keeper) (res []byte, err sdk.Error) { //test this
	attendees := keeper.GetAllAttendees(ctx)
	countAll := len(attendees)

	sort.Slice(attendees, func(i, j int) bool { return attendees[i].Rep > attendees[j].Rep })

	var lb *types.LeaderBoard
	min := types.LeaderBoardCount
	if countAll < types.LeaderBoardCount {
		min = countAll
	}
	top := make([]types.Attendee, min, types.LeaderBoardCount)
	copy(top, attendees)
	for i := range top {
		top[i] = types.Attendee{
			ID:      top[i].ID,
			Address: top[i].Address,
			Name:    top[i].Name,
			Rep:     top[i].Rep,
		}
	}

	lb = types.NewLeaderBoard(countAll, top)
	lb.Time = time.Now()
	res, e := codec.MarshalJSONIndent(keeper.Cdc, lb)
	if e != nil {
		panic("could not marshal result to JSON")
	}

	return
}
