package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/gorilla/mux"
	"net/http"
)

//nolint:gocritic
func bonusGetHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s",
			storeName, keeper.QueryBonus), nil)
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
			storeName, keeper.PrizesKey), nil)
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
		paramType := vars[ScanIDKey]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s",
			storeName, keeper.QueryScans, paramType), nil)
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
		paramType := vars[AddressIDKey]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s",
			storeName, keeper.QueryAttendees, keeper.AddressKey, paramType), nil)
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
		paramType := vars[AttendeeIDKey]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s",
			storeName, keeper.QueryAttendees, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
