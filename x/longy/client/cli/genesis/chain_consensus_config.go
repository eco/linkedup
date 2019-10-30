package genesis

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/types"
	"path/filepath"
	"time"
)

// ConsensusConfigCmd sets the consensus configurations file for the node. Should be run before the chain starts
// up for the first time
func ConsensusConfigCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consensus-config",
		Short: "Set the consensus configurations for the chain. ie. shorten block times",
		RunE: func(_ *cobra.Command, args []string) error {
			setConsensusState(ctx)
			err := setTendermintConsensus(ctx)
			return err
		},
	}

	return cmd
}

func setTendermintConsensus(ctx *server.Context) error {
	genFile := ctx.Config.GenesisFile()
	genDoc, err := types.GenesisDocFromFile(genFile)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	genDoc.ConsensusParams.Block.TimeIotaMs = 1
	err = genutil.ExportGenesisFile(genDoc, genFile)
	return err
}

//setConsensusState sets the consensus configs to quicken block times
func setConsensusState(ctx *server.Context) {
	con := ctx.Config.Consensus
	con.TimeoutPropose = 200 * time.Millisecond
	con.TimeoutProposeDelta = 200 * time.Millisecond
	con.TimeoutPrevote = 200 * time.Millisecond
	con.TimeoutPrevoteDelta = 200 * time.Millisecond
	con.TimeoutPrecommit = 200 * time.Millisecond
	con.TimeoutPrecommitDelta = 200 * time.Millisecond
	con.TimeoutCommit = 800 * time.Millisecond
	cfg.WriteConfigFile(filepath.Join(ctx.Config.RootDir, "config", "config.toml"), ctx.Config)
}
