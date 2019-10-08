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
	ScanIdKey     = "scan_id"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	//longy/attendees/{attendee_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, keeper.QueryAttendees, AttendeeIdKey), attendeeHandler(cliCtx, storeName)).Methods("GET")
	//longy/scans/{scan_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, keeper.QueryScans, ScanIdKey), scanGetHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc("/longy/txs", rest.BroadcastTxRequest(cliCtx)).Methods("POST")
	//r.HandleFunc(fmt.Sprintf("/%s/names", storeName), buyNameHandler(cliCtx)).Methods("POST")
	//r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cliCtx)).Methods("PUT")
	//r.HandleFunc(fmt.Sprintf("/%s/names/{%s}", storeName, restName), resolveNameHandler(cliCtx, storeName)).Methods("GET")
	//r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/whois", storeName, restName), whoIsHandler(cliCtx, storeName)).Methods("GET")
	//r.HandleFunc(fmt.Sprintf("/%s/names", storeName), deleteNameHandler(cliCtx)).Methods("DELETE")
}
