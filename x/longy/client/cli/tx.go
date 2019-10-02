package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/spf13/cobra"
)

//GetTxCmd returns all of the commands to post transaction to the longy module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	longyTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Longy transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return longyTxCmd
}
