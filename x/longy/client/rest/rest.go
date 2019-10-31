package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/eco/longy/x/longy/client/rest/query"
	"github.com/eco/longy/x/longy/internal/querier"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	//LinkedUpHTTPS CORS endpoint for the linked up client
	LinkedUpHTTPS = "https://linkedup.sfbw.io"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
//nolint:gocritic
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	// <storeName>/attendees/{attendee_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, querier.QueryAttendees, query.AttendeeIDKey),
		attendeeHandler(cliCtx, storeName)).Methods(http.MethodGet, http.MethodOptions)

	// <storeName>/attendees/address/{address_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/%s/{%s}", storeName, querier.QueryAttendees, querier.AddressKey,
		query.AddressIDKey), attendeeAddressHandler(cliCtx, storeName)).Methods(http.MethodGet, http.MethodOptions)

	// <storeName>/attendees/{attendee_id}/claimed
	r.HandleFunc(fmt.Sprintf("/%s/attendees/{%s}/claimed", storeName, query.AttendeeIDKey),
		attendeeClaimedHandler(cliCtx, storeName)).Methods(http.MethodGet, http.MethodOptions)

	// <storeName>/attendees/{attendee_id}/keyed
	r.HandleFunc(fmt.Sprintf("/%s/attendees/{%s}/keyed", storeName, query.AttendeeIDKey),
		attendeeKeyedHandler(cliCtx, storeName)).Methods(http.MethodGet, http.MethodOptions)

	// <storeName>/scans/{scan_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", storeName, querier.QueryScans, query.ScanIDKey),
		scanGetHandler(cliCtx, storeName)).Methods(http.MethodGet, http.MethodOptions)

	// <storeName>/prizes
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, querier.PrizesKey),
		prizesGetHandler(cliCtx, storeName)).Methods(http.MethodGet, http.MethodOptions)

	// <storeName>/bonus
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, querier.QueryBonus),
		bonusGetHandler(cliCtx, storeName)).Methods(http.MethodGet, http.MethodOptions)

	// <storeName>/leader
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, querier.LeaderKey),
		query.LeaderBoardHandler(cliCtx, storeName)).Methods(http.MethodGet, http.MethodOptions)

	// <storeName>/winnings?address_id={address_id}
	r.HandleFunc(fmt.Sprintf("/%s/%s", storeName, querier.WinningsKey),
		query.WinningsHandler(cliCtx, storeName)).
		Queries(query.AddressIDKey, fmt.Sprintf("{%s}", query.AddressIDKey)).
		Methods(http.MethodGet, http.MethodOptions)

	// open endpoint to post to in order to claim the prizes of an attendee by passing a sig from the attendee
	r.HandleFunc("/longy/claim", query.ClaimHandler(cliCtx)).Methods(http.MethodPost, http.MethodOptions)

	// open endpoint to post transactions directly to full node
	r.HandleFunc("/longy/txs", rest.BroadcastTxRequest(cliCtx)).Methods(http.MethodPost, http.MethodOptions)

	//  IMPORTANT: you must specify an OPTIONS method matcher for the middleware to set CORS headers
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(CorsMiddleware)
}

//CorsMiddleware adds the CORS header to all the requests
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", LinkedUpHTTPS)
		if r.Method == http.MethodOptions {
			//var headers []string

			//for k := range r.Header {
			//	headers = append(headers, k)
			//	//fmt.Fprintf(w, "Header field %q, Value %q\n", k, v)
			//}
			//allHeaders := strings.Join(headers, ",")
			allHeaders := r.Header.Get("Access-Control-Request-Headers")
			//fmt.Fprintf(w, "Allow headers : %s", allHeaders)
			w.Header().Set("Access-Control-Allow-Headers", allHeaders)
			return
		}
		next.ServeHTTP(w, r)
	})
}
