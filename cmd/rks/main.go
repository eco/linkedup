package main

import (
	"fmt"
	rks "github.com/eco/longy/rekey-service"
	eb "github.com/eco/longy/rekey-service/eventbrite"
	mk "github.com/eco/longy/rekey-service/masterkey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	rootCmd.Flags().Int("port", 1337, "port to bind the rekey service")
	rootCmd.Flags().String("eb-auth-token", "", "eventbrite auth token")
	rootCmd.Flags().Int("eb-event-id", 0, "id associated with the eventbrite event")
	rootCmd.Flags().String("longy-masterkey", "0x0", "master private key for the longy game")
}

var rootCmd = &cobra.Command{
	Use:   "rekey",
	Short: "rekey service for the longest chain game",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		port := viper.GetInt("port")
		authToken := viper.GetString("eb-auth-token")
		eventID := viper.GetInt("eb-event-id")
		key := viper.GetString("longy-masterkey")

		ebSession := eb.CreateSession(authToken, eventID)
		mKey, err := mk.NewMasterKey(key)
		if err != nil {
			return fmt.Errorf("master key: %s", err)
		}

		service := rks.NewService(ebSession, mKey)
		service.StartHTTP(port)

		return nil
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
