package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/spf13/cobra"
)

//GetQueryCmd returns all of the commands to query the longy module
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	longyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the longy module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return longyQueryCmd
}
