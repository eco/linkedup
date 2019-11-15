package query

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/eco/longy/x/longy/crypto"
	"github.com/eco/longy/x/longy/internal/types"
	tmcrypto "github.com/tendermint/tendermint/crypto/secp256k1"
	"net/http"
)

//Claim is the json the claim handler receives
type Claim struct {
	Address string `json:"address"`
	Sig     string `json:"sig"`
}

var signer *crypto.Signer

func init() {
	key := tmcrypto.GenPrivKeySecp256k1([]byte(types.ClaimServiceSeed))
	addr := sdk.AccAddress(key.PubKey().Address())
	signer = crypto.NewSigner(addr, key)
}

//ClaimHandler handles the REST POST to claim the attendee's prizes
//nolint:gocritic
func ClaimHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var claim Claim
		decoder := json.NewDecoder(r.Body)
		//nolint:errcheck
		defer r.Body.Close()
		if err := decoder.Decode(&claim); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid Payload, missing address and sig")
			return
		}

		addr, err := sdk.AccAddressFromBech32(claim.Address)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, sdk.ErrInvalidAddress("invalid address").Error())
			return
		}

		accGetter := auth.NewAccountRetriever(cliCtx)
		if err = accGetter.EnsureExists(addr); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		acc, err := accGetter.GetAccount(addr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = crypto.ValidateSig(acc.GetPubKey(), acc.GetAddress().String(), claim.Sig)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.MsgRedeem{
			Sender:   signer.AccAddress,
			Attendee: addr,
		}

		err = signer.SendTx(&cliCtx, cliCtx.Codec, msg)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, "claim transaction created")
	}
}

//nolint:unparam
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
