package querier

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	// QueryAttendees is the key for attendee gets
	QueryAttendees = "attendees"

	// QueryAttendeeClaimed -
	QueryAttendeeClaimed = "attendee_claimed"

	// QueryAttendeeKeyed -
	QueryAttendeeKeyed = "attendee_keyed"

	// QueryScans is the key for scan gets
	QueryScans = "scans"

	// QueryBonus is the key for checking for the current bonus
	QueryBonus = "bonus"

	// AddressKey is the key for attendee gets by address
	AddressKey = "address"

	// PrizesKey is the key for the event prizes
	PrizesKey = "prizes"

	// LeaderKey is the key for the leader board
	LeaderKey = "leader"

	// WinningsKey is the key for getting the unclaimed prizes of an attendee
	WinningsKey = "winnings"
)

// NewQuerier is the module level router for state queries
//nolint:gocritic,gocyclo
func NewQuerier(keeper keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		queryType := path[0]
		queryArgs := path[1:]

		switch queryType {
		case QueryAttendees:
			if len(queryArgs) > 0 {
				return queryAttendees(ctx, keeper)
			}
			if path[1] == AddressKey {
				queryArgs = path[2:]
				return queryAttendeesByAddr(ctx, queryArgs, keeper)
			}
			return queryAttendee(ctx, queryArgs, keeper)

		case QueryAttendeeClaimed:
			return queryAttendeeClaimed(ctx, queryArgs, keeper)

		case QueryAttendeeKeyed:
			return queryAttendeeKeyed(ctx, queryArgs, keeper)

		case QueryScans:
			if len(queryArgs) > 0 {
				return queryScan(ctx, queryArgs, keeper)
			}

			return queryScans(ctx, keeper)

		case PrizesKey:
			return queryPrizes(ctx, keeper)

		case QueryBonus:
			return queryBonus(ctx, keeper)

		case LeaderKey:
			return leaderBoard(ctx, keeper)

		case WinningsKey:
			return queryWinnings(ctx, keeper, queryArgs)
		}

		return nil, sdk.ErrUnknownRequest("unknown query endpoint")
	}
}
