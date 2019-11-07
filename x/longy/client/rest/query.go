package rest

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/eco/longy/x/longy/client/rest/query"
	"github.com/eco/longy/x/longy/internal/querier"
	longyTypes "github.com/eco/longy/x/longy/internal/types"
	"net/http"

	"github.com/gorilla/mux"
)

//nolint:gocritic
func bonusGetHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s",
			storeName, querier.QueryBonus))
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
func attendeeClaimedHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)[query.AttendeeIDKey]
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s",
			storeName, querier.QueryAttendeeClaimed, id))
		if err != nil {
			if codeType, ok := codeType(err); ok {
				if codeType == longyTypes.AttendeeNotFound {
					http.Error(w, "id not found", http.StatusNotFound)
					return
				}
			}

			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(res) //nolint
	}
}

//nolint:gocritic
func attendeeKeyedHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)[query.AttendeeIDKey]
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s",
			storeName, querier.QueryAttendeeKeyed, id))
		if err != nil {
			if codeType, ok := codeType(err); ok {
				if codeType == longyTypes.AttendeeNotFound {
					http.Error(w, "id not found", http.StatusNotFound)
					return
				}
			}

			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(res) //nolint
	}
}

//nolint:gocritic
func prizesGetHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s",
			storeName, querier.PrizesKey))
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

		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s",
			storeName, querier.QueryScans, paramType))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//nolint:gocritic
func scansGetHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s",
			storeName, querier.QueryScans))
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

		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s/%s",
			storeName, querier.QueryAttendees, querier.AddressKey, paramType))

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

		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s",
			storeName, querier.QueryAttendees, paramType))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//nolint:gocritic
func attendeesHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s",
			storeName, querier.QueryAttendees))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

/** helpers **/
// In tag 0.37.1, the error stringifies into this type. We can extract the code if it's an error of
// this type. We return false if unable
// https://github.com/cosmos/cosmos-sdk/blob/v0.37.1/types/errors.go#L314.
func codeType(err error) (sdk.CodeType, bool) {
	errString := err.Error()
	type codeInBody struct {
		CodeType sdk.CodeType `json:"code"`
	}

	var body codeInBody
	if err := json.Unmarshal([]byte(errString), &body); err != nil {
		return 0, false
	}

	return body.CodeType, true
}
