package main

import (
	"fmt"
	rks "github.com/eco/longy/rekey-service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	rootCmd.Flags().Int("port", 1337, "port to bind the rekey service")
}

var rootCmd = &cobra.Command{
	Use:   "rekey",
	Short: "rekey service for the longest chain game",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())

		port := viper.GetInt("port")
		rks.StartHttpService(port)
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
