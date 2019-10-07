package masterkey

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	longyApp "github.com/eco/longy"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/sirupsen/logrus"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"io"
	"net/http"
	"sync"
	"time"
)

var log = logrus.WithField("module", "masterkey")

var (
	// ErrAlreadyKeyed denotes that this address has already been key'd
	ErrAlreadyKeyed = errors.New("account already key'ed")
)

// MasterKey encapslates the master key for the longy game
type MasterKey struct {
	privKey tmcrypto.PrivKey
	pubKey  tmcrypto.PubKey
	address sdk.AccAddress

	chainID string

	accNum      uint64
	sequenceNum uint64
	seqLock     *sync.Mutex

	cdc         *codec.Codec
	longyCliCtx context.CLIContext
}

// NewMasterKey is the constructor for `Key`. A new secp256k1 is generated if empty.
// The `chainID` is used when generating RekeyTransactions to prevent cross-chain replay attacks
func NewMasterKey(privateKey tmcrypto.PrivKey, restURL, fullNodeURL string, chainID string) (MasterKey, error) {
	cliCtx := context.NewCLIContext().WithNodeURI(fullNodeURL).WithTrustNode(true)
	_, err := cliCtx.Client.Health()
	if err != nil {
		return MasterKey{}, fmt.Errorf("unable to establish connection with the full node")
	}

	// retrieve details about the master account from the rest endpoint
	sdkAddr := sdk.AccAddress(privateKey.PubKey().Address())
	reqURL := restURL + fmt.Sprintf("/auth/accounts/%s", sdkAddr)
	log.Infof("retrieving service key info from %s", reqURL)

	netClient := &http.Client{
		Timeout: 5 * time.Second,
	}
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

		accNum:      account.GetAccountNumber(),
		sequenceNum: account.GetSequence(),
		seqLock:     &sync.Mutex{},

		cdc:         longyApp.MakeCodec(),
		longyCliCtx: cliCtx,
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

	var (
		res sdk.TxResponse

		txBytes []byte
		err     error
	)

	/** Block until we submit the transaction **/
	mk.seqLock.Lock()

	// construct bytes and send to the full node
	txBytes, err = mk.createTxBytes(attendeeAddr, commitment, newPublicKey)
	if err == nil {
		res, err = mk.longyCliCtx.BroadcastTxCommit(txBytes)
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
	}

	mk.seqLock.Unlock()

	return err
}

func (mk *MasterKey) createTxBytes(
	attendeeAddr sdk.AccAddress,
	commitment util.Commitment,
	newPublicKey tmcrypto.PubKey,
) ([]byte, error) {
	msgs := []sdk.Msg{
		longy.NewMsgKey(attendeeAddr, mk.address, newPublicKey, commitment),
	}

	nilFee := auth.NewStdFee(50000, sdk.NewCoins(sdk.NewInt64Coin("longy", 0)))
	signBytes := auth.StdSignBytes(mk.chainID, mk.accNum, mk.sequenceNum, nilFee, msgs, "")

	// sign the message with the master private key
	sig, err := mk.privKey.Sign(signBytes)
	if err != nil {
		return nil, err
	}
	stdSig := auth.StdSignature{PubKey: mk.pubKey, Signature: sig}
	tx := auth.NewStdTx(msgs, nilFee, []auth.StdSignature{stdSig}, "")

	return auth.DefaultTxEncoder(mk.cdc)(tx)
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
