package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"os"
	"sync"
)

//Signer is the signer for posting redeem tx
type Signer struct {
	AccAddress sdk.AccAddress
	PrivKey    secp256k1.PrivKeySecp256k1
	seqLock    *sync.Mutex
}

//NewSigner is the constructor for signer
func NewSigner(accAddress sdk.AccAddress, privKey secp256k1.PrivKeySecp256k1) *Signer {
	return &Signer{
		AccAddress: accAddress,
		PrivKey:    privKey,
		seqLock:    &sync.Mutex{},
	}
}

// SendTx calls the test handler with a message.
func (s *Signer) SendTx(cliContext *context.CLIContext, cdc *codec.Codec, msg sdk.Msg) (res sdk.Result) {
	s.seqLock.Lock()
	txBldr := s.senderTxContext(cliContext, cdc)
	err := s.completeAndBroadcastTxCLI(*cliContext, txBldr, []sdk.Msg{msg})
	s.seqLock.Unlock()
	if err != nil {
		//res = parseResponse(err.Error())
		fmt.Println(err)
	}

	return
}

// SenderTxContext creates a new TxBuilder
func (s *Signer) senderTxContext(cliContext *context.CLIContext, cdc *codec.Codec) auth.TxBuilder {
	cliContext.FromAddress = s.AccAddress
	cliContext.BroadcastMode = client.BroadcastSync
	txEncoder := utils.GetTxEncoder(cdc)
	return auth.NewTxBuilderFromCLI().WithTxEncoder(txEncoder)
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

	if txBldr.SimulateAndExecute() || cliContext.Simulate {
		txBldr, err = utils.EnrichWithGas(txBldr, cliContext, msgs)
		if err != nil {
			return err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if cliContext.Simulate {
		return nil
	}

	// build and sign the transaction
	txBytes, err := buildAndSign(txBldr, s.PrivKey, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := cliContext.BroadcastTxCommit(txBytes)
	if err != nil {
		return err
	}

	err = cliContext.PrintOutput(res)
	return err
}

//nolint[hugeParam]
// Function signature taken from the Cosmos SDK.
func buildAndSign(bldr auth.TxBuilder, priv crypto.PrivKey, msgs []sdk.Msg) ([]byte, error) {
	msg, err := bldr.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	sig, err := makeSignature(priv, msg)
	if err != nil {
		return nil, err
	}

	return bldr.TxEncoder()(auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo))
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
