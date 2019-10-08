package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/gorilla/mux"
)

const (
	AttendeeIdKey = "attendee_id"
	AddressIdKey  = "address_id"
	ScanIdKey     = "scan_id"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	//longy/attendees/{attendee_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, keeper.QueryAttendees, AttendeeIdKey), attendeeHandler(cliCtx, storeName)).Methods("GET")
	//longy/attendees/address/{address_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/%s/{%s}", storeName, keeper.QueryAttendees, keeper.AddressKey, AddressIdKey), attendeeAddressHandler(cliCtx, storeName)).Methods("GET")
	//longy/scans/{scan_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, keeper.QueryScans, ScanIdKey), scanGetHandler(cliCtx, storeName)).Methods("GET")

	//open endpoint to post transactions directly to full node
	r.HandleFunc("/longy/txs", rest.BroadcastTxRequest(cliCtx)).Methods("POST")
}
