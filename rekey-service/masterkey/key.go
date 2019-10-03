package masterkey

import (
	"encoding/hex"
	"encoding/json"
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
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"io"
	"net/http"
	"sync"
)

const (
	secp256k1PrivKeyLen = 32
)

var log = logrus.WithField("module", "masterkey")

// Key encapslates the master key for the
// longey game
type Key struct {
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
func NewMasterKey(privateKey tmcrypto.PrivKey, restURL, fullNodeURL string, chainID string) (Key, error) {
	cliCtx := context.NewCLIContext().
		WithNodeURI(fullNodeURL).WithTrustNode(true) // current release doesn't support `WithChainID`
	_, err := cliCtx.Client.Health()
	if err != nil {
		return Key{}, fmt.Errorf("unable to establish connection with the full node")
	}

	// retrieve details about the master account from the rest endpoint
	reqURL := restURL + "/auth/accounts/{addr goes here}"
	resp, err := http.Get(reqURL)
	acc, err := parseAccountFromBody(resp.Body)
	if err != nil {
		return Key{}, fmt.Errorf("unable to request master account information from the full node")
	}

	k := Key{
		privKey: privateKey,
		pubKey:  privateKey.PubKey(),

		accNum:      acc.GetAccountNumber(),
		sequenceNum: acc.GetSequence(),
		seqLock:     &sync.Mutex{},

		cdc:         longyApp.MakeCodec(),
		longyCliCtx: cliCtx,
	}

	return k, nil
}

// Secp256k1FromHex parses the hex-encoded `key` string
func Secp256k1FromHex(key string) (tmcrypto.PrivKey, error) {
	if len(key) == 0 {
		log.Info("provided key is empty. generating a new Secp256k1 key")
		return secp256k1.GenPrivKey(), nil
	}

	bytes, err := hex.DecodeString(util.TrimHex(key))
	if err != nil {
		return nil, fmt.Errorf("hex decoding: %s", err)
	} else if len(bytes) != secp256k1PrivKeyLen {
		return nil, fmt.Errorf("invalid key byte length. expected: %d, got: %d",
			secp256k1PrivKeyLen, len(bytes))
	}

	var privateKey [secp256k1PrivKeyLen]byte
	copied := copy(privateKey[:], bytes)
	if copied != secp256k1PrivKeyLen {
		errMsg := fmt.Sprintf("expected to copy over %d bytes into the secp256k1 private key",
			secp256k1PrivKeyLen)
		panic(errMsg)
	}

	return secp256k1.PrivKeySecp256k1(privateKey), nil
}

// SendRekeyTransaction generates a `RekeyMsg`, authorized by the master key. The transaction bytes
// generated are created using the cosmos-sdk/x/auth module's StdSignDoc.
func (mk Key) SendRekeyTransaction(attendeeID string, secret []byte, newPublicKey tmcrypto.PubKey) error {
	var err error

	/** Block until we submit the transaction **/
	mk.seqLock.Lock()

	// construct bytes and send to the full node
	txBytes, err := mk.createTxBytes(attendeeID, secret, newPublicKey)
	if err == nil {
		_, err = mk.longyCliCtx.BroadcastTxSync(txBytes)
	}

	mk.seqLock.Unlock()

	return err
}

func (mk Key) createTxBytes(attendeeID string, secret []byte, newPublicKey tmcrypto.PubKey) ([]byte, error) {
	attendeeAddr := util.IDToAddress(attendeeID)
	msgs := []sdk.Msg{longy.NewRekeyMsg(attendeeAddr, mk.address, newPublicKey, secret)}
	signBytes := auth.StdSignBytes(
		mk.chainID,
		mk.accNum,
		mk.sequenceNum,
		nil,
		msgs,
		"",
	)

	// sign with the private key
	sig, err := mk.privKey.Sign(signBytes)
	if err != nil {
		return nil, err
	}
	stdSig := auth.StdSignature{
		PubKey:    mk.pubKey,
		Signature: sig,
	}
	tx := auth.NewStdTx(msgs, nil, []auth.StdSignature{stdSig}, "")

	return auth.DefaultTxEncoder(mk.cdc)(tx)
}

func parseAccountFromBody(body io.ReadCloser) (auth.Account, error) {
	defer body.Close()
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
	err = json.Unmarshal(accBody, &acc)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}
