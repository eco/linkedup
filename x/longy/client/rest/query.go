package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
	"net/http"

	"github.com/gorilla/mux"
)

var signer *keeper.Signer

func init() {
	key := tmcrypto.GenPrivKeySecp256k1([]byte("master"))
	addr := sdk.AccAddress(key.PubKey().Address())
	signer = keeper.NewSigner(addr, key)
}

//nolint:gocritic
func redeemHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addressID := vars[AddressIDKey]
		sig := vars[SigKey]

		res, _, _ := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s",
			storeName, keeper.RedeemKey, addressID, sig), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)

		addr, _ := sdk.AccAddressFromBech32(addressID) //err checked in query

		msg := types.MsgRedeem{
			Sender:   signer.AccAddress,
			Attendee: addr,
		}
		signer.SendTx(&cliCtx, cliCtx.Codec, msg)
	}
}

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
