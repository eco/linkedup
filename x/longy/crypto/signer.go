package crypto

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/tendermint/tendermint/crypto"
	"sync"
)

//Signer is the signer for posting redeem tx
type Signer struct {
	AccAddress sdk.AccAddress
	PrivKey    crypto.PrivKey
	seqLock    *sync.Mutex
}

//NewSigner is the constructor for signer
func NewSigner(accAddress sdk.AccAddress, privKey crypto.PrivKey) *Signer {
	return &Signer{
		AccAddress: accAddress,
		PrivKey:    privKey,
		seqLock:    &sync.Mutex{},
	}
}

// SendTx calls the test handler with a message.
func (s *Signer) SendTx(cliContext *context.CLIContext, cdc *codec.Codec, msg sdk.Msg) error {
	s.seqLock.Lock()
	txBldr, err := s.senderTxContext(cliContext, cdc)
	if err != nil {
		//res = parseResponse(err.Error())
		//fmt.Println(err)
		return err
	}
	err = s.completeAndBroadcastTxCLI(*cliContext, txBldr, []sdk.Msg{msg})
	s.seqLock.Unlock()
	if err != nil {
		//res = parseResponse(err.Error())
		//fmt.Println(err)
		return err
	}

	return nil
}

// SenderTxContext creates a new TxBuilder
func (s *Signer) senderTxContext(cliContext *context.CLIContext, cdc *codec.Codec) (auth.TxBuilder, error) {
	cliContext.FromAddress = s.AccAddress
	//cliContext.BroadcastMode = client.BroadcastSync//todo switch back to this setting
	cliContext.BroadcastMode = client.BroadcastBlock
	txEncoder := utils.GetTxEncoder(cdc)
	txBuilder := auth.NewTxBuilderFromCLI().WithTxEncoder(txEncoder)

	num, seq, err := auth.NewAccountRetriever(cliContext).GetAccountNumberSequence(s.AccAddress)
	if err != nil {
		return txBuilder, err
	}

	return txBuilder.WithAccountNumber(num).WithSequence(seq), nil
}

// nolint[hugeParam]
// Function signature taken from the Cosmos SDK.
// CompleteAndBroadcastTxCLI broadcasts a Tx.
func (s *Signer) completeAndBroadcastTxCLI(cliContext context.CLIContext, txBldr auth.TxBuilder,
	msgs []sdk.Msg) error {

	txBldr, err := utils.PrepareTxBuilder(txBldr, cliContext)
	if err != nil {
		return err
	}

	// build and sign the transaction
	txBytes, err := s.buildAndSign(cliContext.Codec, txBldr, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := cliContext.BroadcastTxCommit(txBytes)
	if err != nil {
		return err
	}

	if sdk.CodeType(res.Code) != sdk.CodeOK {
		return fmt.Errorf(res.RawLog)
	}

	err = cliContext.PrintOutput(res)
	return err
}

//nolint[hugeParam]
// Function signature taken from the Cosmos SDK.
func (s *Signer) buildAndSign(cdc *codec.Codec, bldr auth.TxBuilder, msgs []sdk.Msg) ([]byte, error) {
	msg, err := bldr.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(msg)
	fmt.Println(string(b))
	sig, err := makeSignature(s.PrivKey, msg)
	if err != nil {
		return nil, err
	}

	fee := msg.Fee
	fee.Gas = fee.Gas * 2

	return bldr.TxEncoder()(auth.NewStdTx(msg.Msgs, fee, []auth.StdSignature{sig}, msg.Memo))
}

//nolint[hugeParam]
// Function signature taken from the Cosmos SDK.
func makeSignature(priv crypto.PrivKey, msg auth.StdSignMsg) (sig auth.StdSignature, err error) {

	sigBytes, err := priv.Sign(msg.Bytes())
	if err != nil {
		return
	}

	pubkey := priv.PubKey()

	return auth.StdSignature{
		PubKey:    pubkey,
		Signature: sigBytes,
	}, nil
}
