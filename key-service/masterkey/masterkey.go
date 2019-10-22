package masterkey

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	longyApp "github.com/eco/longy"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/sirupsen/logrus"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var log = logrus.WithField("module", "masterkey")

var (
	netClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	// ErrAlreadyKeyed denotes that this address has already been key'd
	ErrAlreadyKeyed = errors.New("account already key'ed")

	// ErrTxSubmission
	ErrTxSubmission = errors.New("unable to complete tx submission")
)

// MasterKey encapslates the master key for the longy game
type MasterKey struct {
	privKey tmcrypto.PrivKey
	pubKey  tmcrypto.PubKey
	address sdk.AccAddress

	chainID string
	restURL string

	accNum      uint64
	sequenceNum uint64
	seqLock     *sync.Mutex

	cdc *codec.Codec
}

// NewMasterKey is the constructor for `Key`. A new secp256k1 is generated if empty.
// The `chainID` is used when generating RekeyTransactions to prevent cross-chain replay attacks
func NewMasterKey(privateKey tmcrypto.PrivKey, restURL, chainID string) (MasterKey, error) {

	// retrieve details about the master account from the rest endpoint
	sdkAddr := sdk.AccAddress(privateKey.PubKey().Address())
	reqURL := restURL + fmt.Sprintf("/auth/accounts/%s", sdkAddr)
	log.Infof("retrieving service key info from %s", reqURL)

	resp, err := netClient.Get(reqURL)
	if err != nil {
		return MasterKey{}, fmt.Errorf("unable to establish connection to the rest service: %s", err)
	}
	account, err := parseAccountFromBody(resp.Body)
	if err != nil {
		return MasterKey{}, fmt.Errorf("unable to request master account information from the full node: %s", err)
	}

	k := MasterKey{
		privKey: privateKey,
		pubKey:  privateKey.PubKey(),
		address: sdkAddr,

		chainID: chainID,
		restURL: restURL,

		accNum:      account.GetAccountNumber(),
		sequenceNum: account.GetSequence(),
		seqLock:     &sync.Mutex{},

		cdc: longyApp.MakeCodec(),
	}

	_ = resp.Body.Close()
	log.Infof("constructed master key. Chain-Id=%s, AccountNum=%d, SequenceNum=%d", k.chainID, k.accNum, k.sequenceNum)
	return k, nil
}

// SendKeyTransaction generates a `RekeyMsg`, authorized by the master key. The transaction bytes
// generated are created using the cosmos-sdk/x/auth module's StdSignDoc.
func (mk *MasterKey) SendKeyTransaction(
	attendeeAddr sdk.AccAddress,
	newPublicKey tmcrypto.PubKey,
	commitment util.Commitment,
) error {

	/** Block until we submit the transaction **/
	mk.seqLock.Lock()

	// create and broadcast the transaction
	keyMsg := longy.NewMsgKey(attendeeAddr, mk.address, newPublicKey, commitment)
	tx, err := mk.createKeyTx(keyMsg)
	res, err := mk.broadcastTx(*tx)
	if err != nil { // nolint
		log.WithError(err).Info("failed transaction submission")
	} else {
		if res.Code != 0 {
			if res.Code == uint32(longy.CodeAttendeeKeyed) {
				err = ErrAlreadyKeyed
			} else {
				log.WithField("raw_log", res.RawLog).Info("failed tx response")
				err = fmt.Errorf("failed tx")
			}
		}

		mk.sequenceNum++
	}

	mk.seqLock.Unlock()

	return err
}

//nolint
func (mk *MasterKey) broadcastTx(tx auth.StdTx) (*sdk.TxResponse, error) {
	reqURL := mk.restURL + "/longy/txs"

	body := struct {
		Tx   auth.StdTx `json:"tx"`
		Mode string     `json:"mode"`
	}{Tx: tx, Mode: "block"}

	bz, err := mk.cdc.MarshalJSON(body)
	if err != nil {
		log.WithError(err).Error("marshal tx body")
		return nil, err
	}

	resp, err := netClient.Post(reqURL, "application/json", bytes.NewReader(bz))
	if err != nil {
		log.WithError(err).Error("http tx submission")
		return nil, err
	} else if resp.Status == "500" {
		return nil, ErrTxSubmission
	}
	bz, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("reading from tx submission resp body")
		panic(err)
	}

	var res sdk.TxResponse
	err = mk.cdc.UnmarshalJSON(bz, &res)
	if err != nil {
		// key-service is broken at this point. What is the right sequence number?
		panic(err)
	}

	return &res, nil
}

//nolint
func (mk *MasterKey) createKeyTx(keyMsg longy.MsgKey) (*auth.StdTx, error) {
	msgs := []sdk.Msg{keyMsg}

	nilFee := auth.NewStdFee(50000, sdk.NewCoins(sdk.NewInt64Coin("longy", 0)))
	signBytes := auth.StdSignBytes(mk.chainID, mk.accNum, mk.sequenceNum, nilFee, msgs, "")

	// sign the message with the master private key
	sig, err := mk.privKey.Sign(signBytes)
	if err != nil {
		return nil, err
	}
	stdSig := auth.StdSignature{PubKey: mk.pubKey, Signature: sig}
	tx := auth.NewStdTx(msgs, nilFee, []auth.StdSignature{stdSig}, "")

	return &tx, nil
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
