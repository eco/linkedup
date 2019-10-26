package longyclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	longyCfg "github.com/eco/longy/key-service/config"
)

var netClient = &http.Client{
	Timeout: 10 * time.Second,
}

// GetAccount -
func GetAccount(addr sdk.AccAddress) (auth.Account, error) {
	restURL := longyCfg.LongyRestURL()
	reqURL := restURL + fmt.Sprintf("/auth/accounts/%s", addr)
	resp, err := netClient.Get(reqURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return parseAccountFromBody(resp.Body)
}

// BroadcastAuthTx -
func BroadcastAuthTx(tx auth.StdTx) (*sdk.TxResponse, error) {
	restURL := longyCfg.LongyRestURL()
	reqURL := restURL + "/longy/txs"

	body := struct {
		Tx   auth.StdTx `json:"tx"`
		Mode string     `json:"mode"`
	}{Tx: tx, Mode: "block"}

	bz, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := netClient.Post(reqURL, "application/json", bytes.NewReader(bz))
	if err != nil {
		return nil, err
	}

	// If we error in the following statements, key-service is broken at this point
	// What is the right sequence number? It could have incremented, we cannot know forsure
	//
	// Right course of action is to panic and restart the service so the master account
	// re-syncs with the full node

	bz, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var res sdk.TxResponse
	if err = json.Unmarshal(bz, &res); err != nil {
		panic(err)
	}

	return &res, nil
}

func parseAccountFromBody(body io.ReadCloser) (auth.Account, error) {
	decoder := json.NewDecoder(body)

	var b map[string]json.RawMessage
	err := decoder.Decode(&b)
	if err != nil {
		return nil, err
	}

	var acc auth.BaseAccount
	accBody, ok := b["result"]
	if !ok {
		return nil, fmt.Errorf("result not present in the body")
	}
	err = auth.ModuleCdc.UnmarshalJSON(accBody, &acc)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}
