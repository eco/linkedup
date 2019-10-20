package querier

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	//QueryAttendees  is the key for attendee gets
	QueryAttendees = "attendees"
	//QueryScans is the key for scan gets
	QueryScans = "scans"
	//QueryBonus is the key for checking for the current bonus
	QueryBonus = "bonus"
	//AddressKey is the key for attendee gets by address
	AddressKey = "address"
	//PrizesKey is the key for the event prizes
	PrizesKey = "prizes"
	//LeaderKey is the key for the leader board
	LeaderKey = "leader"
	//RedeemKey is the key for the redeem event
	RedeemKey = "redeem"
)

// NewQuerier is the module level router for state queries
//nolint:gocritic
func NewQuerier(keeper keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		queryType := path[0]
		queryArgs := path[1:]

		switch queryType {
		case QueryAttendees:
			if path[1] == AddressKey {
				queryArgs = path[2:]
				return queryAttendeesByAddr(ctx, queryArgs, keeper)
			}
			return queryAttendees(ctx, queryArgs, keeper)
		case QueryScans:
			return queryScans(ctx, queryArgs, keeper)
		case PrizesKey:
			return queryPrizes(ctx, keeper)
		case QueryBonus:
			return queryBonus(ctx, keeper)
		case LeaderKey:
			return leaderBoard(ctx, keeper)
		case RedeemKey:
			return queryRedeem(ctx, keeper, queryArgs)
		default:
			break
		}

		return nil, sdk.ErrUnknownRequest("unknown query endpoint")
	}
}
