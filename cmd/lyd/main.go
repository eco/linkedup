package main

import (
	"github.com/cosmos/cosmos-sdk/server"
	app "github.com/eco/longy"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"io"
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "lyd",
		Short:             "longy the longest chain validator daemon",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, newApp, nil)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "LY", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewLongyApp(logger, db)
}
