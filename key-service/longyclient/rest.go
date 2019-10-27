package longyclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	longyApp "github.com/eco/longy"
	longyCfg "github.com/eco/longy/key-service/config"
	"github.com/sirupsen/logrus"
)

var (
	netClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	cdc = longyApp.MakeCodec()

	log = logrus.WithField("module", "longyclient")
)

// IsAttendeeKeyed -
func IsAttendeeKeyed(id int) (bool, error) {
	if id < 0 {
		return false, fmt.Errorf("id must be a positive integer")
	}

	restURL := longyCfg.LongyRestURL()
	reqURL := restURL + fmt.Sprintf("/longy/attendees/%d/keyed", id)
	resp, err := netClient.Get(reqURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	keyed, err := strconv.ParseBool(string(body))
	if err != nil {
		return false, fmt.Errorf("unexpected attendee keyed response, %s", err)
	}

	return keyed, nil
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

// BroadcastAuthTx - mode = "sync|async|block"
func BroadcastAuthTx(tx auth.StdTx, mode string) (*sdk.TxResponse, error) {
	if mode != "sync" && mode != "async" && mode != "block" {
		return nil, fmt.Errorf("incorrect broadcast mode")
	}
	restURL := longyCfg.LongyRestURL()
	reqURL := restURL + "/longy/txs"
	body := struct {
		Tx   auth.StdTx `json:"tx"`
		Mode string     `json:"mode"`
	}{Tx: tx, Mode: mode}

	bz, err := cdc.MarshalJSON(body)
	if err != nil {
		return nil, err
	}

	resp, err := netClient.Post(reqURL, "application/json", bytes.NewReader(bz))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If we error in the following statements, key-service is broken at this point
	// What is the right sequence number? It could have incremented, we cannot know forsure
	//
	// Right course of action is to panic and restart the service so the master account
	// re-syncs with the full node

	bz, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	} else if resp.StatusCode != 200 {
		log.WithField("status", resp.Status).WithField("body", string(bz)).
			Error("non-ok broadcast tx submission response")

		return nil, fmt.Errorf("broadcast tx response")
	}

	var res sdk.TxResponse
	if err = cdc.UnmarshalJSON(bz, &res); err != nil {
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
