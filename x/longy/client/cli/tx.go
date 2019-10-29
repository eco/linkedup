package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	longyCfg "github.com/eco/longy/key-service/config"
	longyClnt "github.com/eco/longy/key-service/longyclient"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

//nolint
func init() {
	// default value uses "bonus" as the seed
	longyTxCmd.PersistentFlags().String("private-key", "5f220b4a45832c0710e50d99fb0202a8554dcfe4e8593939d2820a689cbff212",
		"hex-encoded secp256k1 private key of the bonus service account")
	longyTxCmd.PersistentFlags().String("longy-rest-url", "http://localhost:1317", "scheme://host:port of the longy rest service")
	viper.BindPFlag("private-key", longyTxCmd.PersistentFlags().Lookup("private-key"))
	viper.BindPFlag("longy-rest-url", longyTxCmd.PersistentFlags().Lookup("longy-rest-url"))
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
			chainID := viper.GetString(client.FlagChainID)
			restURL := viper.GetString("longy-rest-url")
			fmt.Printf("Longy Rest Service URL: %s\n", restURL)
			longyCfg.SetLongyRestURL(restURL)

			bonusAmt, err := strconv.Atoi(args[0])
			if err != nil || bonusAmt <= 0 {
				return fmt.Errorf("multiplier must be a positive number > 0 in decimal format")
			}

			/** read in the private key and retrieve the bonus account */
			bonusAccount, privKey, err := readBonusAccountFromViper()
			if err != nil {
				return fmt.Errorf("master account: %s", err)
			}

			/** construct the message **/
			bonusMsg := longy.NewMsgBonus(uint(bonusAmt), bonusAccount.GetAddress())

			/** send the message **/
			tx, err := createAuthTx(
				cdc, bonusMsg, privKey, chainID,
				bonusAccount.GetAccountNumber(), bonusAccount.GetSequence())
			if err != nil {
				return fmt.Errorf("tx creation: %s", err)
			}

			res, err := longyClnt.BroadcastAuthTx(tx, "block")
			if err != nil {
				return fmt.Errorf("tx sbumission: %s", err)
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
			chainID := viper.GetString(client.FlagChainID)
			restURL := viper.GetString("longy-rest-url")
			fmt.Printf("Longy Rest Service URL: %s\n", restURL)
			longyCfg.SetLongyRestURL(restURL)

			/** read in the private key and retrieve the master account */
			bonusAccount, privKey, err := readBonusAccountFromViper()
			if err != nil {
				return fmt.Errorf("master account: %s", err)
			}

			/** construct the message **/
			clearBonusMsg := longy.NewMsgClearBonus(bonusAccount.GetAddress())

			/** send the message **/
			tx, err := createAuthTx(
				cdc, clearBonusMsg, privKey, chainID,
				bonusAccount.GetAccountNumber(), bonusAccount.GetSequence())
			if err != nil {
				return fmt.Errorf("tx creation: %s", err)
			}

			res, err := longyClnt.BroadcastAuthTx(tx, "block")
			if err != nil {
				return fmt.Errorf("tx submission: %s", err)
			}

			fmt.Printf("Transaction Response:\n%v\n", res)

			return nil
		},
	}
}

func readBonusAccountFromViper() (auth.Account, tmcrypto.PrivKey, error) {
	privKey, err := util.Secp256k1FromHex(viper.GetString("private-key"))
	if err != nil {
		return nil, nil, fmt.Errorf("private key: %s\n", err)
	}

	// retrieve the bonus account
	addr := sdk.AccAddress(privKey.PubKey().Address())
	bonusAccount, err := longyClnt.GetAccount(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("retrieving bonus account: %s\n", err)
	}

	return bonusAccount, privKey, nil
}

func createAuthTx(
	cdc *codec.Codec,
	msg sdk.Msg,
	privKey tmcrypto.PrivKey,
	chainID string,
	accNum, seqNum uint64) (*auth.StdTx, error) {

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

	return &tx, nil
}
