package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"strconv"
)

func init() {
	longyTxCmd.PersistentFlags().String("private-key", "", "hex-encoded secp256k1 private key of the master account") //nolint
	viper.BindPFlag("private-key", longyTxCmd.PersistentFlags().Lookup("private-key"))                                //nolint
}

var longyTxCmd = &cobra.Command{
	Use:   types.ModuleName,
	Short: "Longy transaction subcommands",
	RunE:  client.ValidateCmd,
}

//GetTxCmd returns all of the commands to post transaction to the longy module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	longyTxCmd.AddCommand(
		createBonusCmd(cdc),
		clearBonusCmd(cdc),
	)

	return longyTxCmd
}

func createBonusCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-bonus <multiplier>",
		Short: "create a bonus for sponsor scans",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags()) //nolint
			cliCtx := context.NewCLIContext().WithTrustNode(true).WithCodec(cdc)
			chainID := viper.GetString(client.FlagChainID)

			bonusAmt, err := strconv.Atoi(args[0])
			if err != nil || bonusAmt <= 0 {
				return fmt.Errorf("multiplier must be a positive number > 0 in decimal format")
			}

			/** read in the private key and retrieve the master account */
			masterAcc, privKey, err := getMasterAccountFromViper(&cliCtx)
			if err != nil {
				return fmt.Errorf("master account: %s", err)
			}

			/** construct the message **/
			bonusMsg := longy.NewMsgBonus(uint(bonusAmt), masterAcc.GetAddress())

			/** send the message **/
			txBytes, err := createTxBytes(
				cdc, bonusMsg, privKey, chainID,
				masterAcc.GetAccountNumber(), masterAcc.GetSequence())
			if err != nil {
				return fmt.Errorf("tx bytes: %s", err)
			}

			res, err := cliCtx.BroadcastTxCommit(txBytes)
			if err != nil {
				return fmt.Errorf("tx submission: %s", err)
			}

			fmt.Printf("Transaction Response:\n%v\n", res)

			return nil
		},
	}
}

func clearBonusCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "clear-bonus",
		Short: "clear the bonus period",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags()) //nolint
			cliCtx := context.NewCLIContext().WithTrustNode(true).WithCodec(cdc)
			chainID := viper.GetString(client.FlagChainID)

			/** read in the private key and retrieve the master account */
			masterAcc, privKey, err := getMasterAccountFromViper(&cliCtx)
			if err != nil {
				return fmt.Errorf("master account: %s", err)
			}

			/** construct the message **/
			clearBonusMsg := longy.NewMsgClearBonus(masterAcc.GetAddress())

			/** send the message **/
			txBytes, err := createTxBytes(
				cdc, clearBonusMsg, privKey, chainID,
				masterAcc.GetAccountNumber(), masterAcc.GetSequence())
			if err != nil {
				return fmt.Errorf("tx bytes: %s", err)
			}

			res, err := cliCtx.BroadcastTxCommit(txBytes)
			if err != nil {
				return fmt.Errorf("tx submission: %s", err)
			}

			fmt.Printf("Transaction Response:\n%v\n", res)

			return nil
		},
	}
}

//nolint
func getMasterAccountFromViper(cliCtx *context.CLIContext) (auth.Account, tmcrypto.PrivKey, error) {
	accRetriever := auth.NewAccountRetriever(cliCtx)

	keyStr := viper.GetString("private-key")
	if len(keyStr) == 0 {
		return nil, nil, fmt.Errorf("empty private key")
	}
	privKey, err := util.Secp256k1FromHex(keyStr)
	if err != nil {
		return nil, nil, fmt.Errorf("private-key: %s", err)
	}
	masterAddr := sdk.AccAddress(privKey.PubKey().Address())

	acc, err := accRetriever.GetAccount(masterAddr)
	if err != nil {
		return nil, nil, err
	}

	return acc, privKey, nil
}

func createTxBytes(
	cdc *codec.Codec,
	msg sdk.Msg,
	privKey tmcrypto.PrivKey,
	chainID string,
	accNum, seqNum uint64) ([]byte, error) {

	msgs := []sdk.Msg{msg}
	nilFee := auth.NewStdFee(50000, sdk.NewCoins(sdk.NewInt64Coin("longy", 0)))
	signBytes := auth.StdSignBytes(chainID, accNum, seqNum, nilFee, msgs, "")

	// sign the message with the master private key
	sig, err := privKey.Sign(signBytes)
	if err != nil {
		return nil, err
	}
	stdSig := auth.StdSignature{PubKey: privKey.PubKey(), Signature: sig}
	tx := auth.NewStdTx(msgs, nilFee, []auth.StdSignature{stdSig}, "")

	return auth.DefaultTxEncoder(cdc)(tx)
}
