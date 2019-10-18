package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/gorilla/mux"
)

const (
	//AttendeeIDKey is the attribute key for attendee id
	AttendeeIDKey = "attendee_id"
	//AddressIDKey is the attribute key for attendee address
	AddressIDKey = "address_id"
	//ScanIDKey is the attribute key for scan id
	ScanIDKey = "scan_id"
	//BadgeIDKey is the attribute key for badge id
	BadgeIDKey = "badge_id"
	//SigKey is the attribute key for the sig
	SigKey = "sig"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
//nolint:gocritic
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	//longy/attendees/{attendee_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, keeper.QueryAttendees, AttendeeIDKey),
		attendeeHandler(cliCtx, storeName)).Methods("GET")

	//longy/attendees/address/{address_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/%s/{%s}", storeName, keeper.QueryAttendees, keeper.AddressKey,
		AddressIDKey), attendeeAddressHandler(cliCtx, storeName)).Methods("GET")

	//longy/scans/{scan_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, keeper.QueryScans, ScanIDKey),
		scanGetHandler(cliCtx, storeName)).Methods("GET")

	//longy/prizes
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, keeper.PrizesKey),
		prizesGetHandler(cliCtx, storeName)).Methods("GET")

	//longy/bonus
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, keeper.QueryBonus),
		bonusGetHandler(cliCtx, storeName)).Methods("GET")

	//longy/redeem?address_id={address_id}&sig={sig}
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, keeper.RedeemKey),
		redeemHandler(cliCtx, storeName)).
		Queries(AddressIDKey, fmt.Sprintf("{%s}", AddressIDKey), SigKey, fmt.Sprintf("{%s}", SigKey)).
		Methods("GET")

	//open endpoint to post transactions directly to full node
	r.HandleFunc("/longy/txs", rest.BroadcastTxRequest(cliCtx)).Methods("POST")
}
