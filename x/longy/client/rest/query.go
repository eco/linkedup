package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/eco/longy/x/longy/client/rest/query"
	"github.com/eco/longy/x/longy/internal/querier"
	"net/http"

	"github.com/gorilla/mux"
)

//nolint:gocritic
func bonusGetHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
			storeName, querier.QueryBonus), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		} else if res == nil {
			// no bonus found
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(res) //nolint
	}
}

//nolint:gocritic
func prizesGetHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
			storeName, querier.PrizesKey), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//nolint:gocritic
func scanGetHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[query.ScanIDKey]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s",
			storeName, querier.QueryScans, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//nolint:gocritic
func attendeeAddressHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[query.AddressIDKey]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s",
			storeName, querier.QueryAttendees, querier.AddressKey, paramType), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//nolint:gocritic
func attendeeHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[query.AttendeeIDKey]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s",
			storeName, querier.QueryAttendees, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
