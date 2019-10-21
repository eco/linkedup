package query

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	//"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/querier"
	//"github.com/eco/longy/x/longy/internal/types"
	//tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
	"net/http"

	"github.com/gorilla/mux"
)

//var signer *keeper.Signer
//
//func init() {
//	key := tmcrypto.GenPrivKeySecp256k1([]byte("master"))
//	addr := sdk.AccAddress(key.PubKey().Address())
//	signer = keeper.NewSigner(addr, key)
//}

//RedeemHandler responds to redeem prize queries for which prizes haven't been claimed yet for an attendee
//nolint:gocritic
func RedeemHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addressID := vars[AddressIDKey]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s",
			storeName, querier.RedeemKey, addressID), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)

		//addr, _ := sdk.AccAddressFromBech32(addressID) //err checked in query
		//
		//msg := types.MsgRedeem{
		//	Sender:   signer.AccAddress,
		//	Attendee: addr,
		//}
		//signer.SendTx(&cliCtx, cliCtx.Codec, msg)
	}
}
