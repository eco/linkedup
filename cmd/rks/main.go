package main

import (
	"fmt"
	rks "github.com/eco/longy/rekey-service"
	eb "github.com/eco/longy/rekey-service/eventbrite"
	"github.com/eco/longy/rekey-service/mail"
	mk "github.com/eco/longy/rekey-service/masterkey"
	"github.com/eco/longy/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	// TODO: remote these defaults before making the repo public
	rootCmd.Flags().Int("port", 1337, "port to bind the rekey service")
	rootCmd.Flags().String("longy-masterkey", "", "master private key for the longy game. A new keypair will be generated if not provided")
	rootCmd.Flags().String("smtp-server", "smtp.gmail.com:587", "host:port of the smtp server")
	rootCmd.Flags().String("smtp-username", "testecolongy@gmail.com", "username of the email account")
	rootCmd.Flags().String("smtp-password", "2019longygame", "password of the email account")
	rootCmd.Flags().String("eb-auth-token", "", "eventbrite authorization token")
	rootCmd.Flags().Int("eb-event-id", 0, "id associated with the eventbrite event")
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
		smtpServer := viper.GetString("smtp-server")
		smtpUsername := viper.GetString("smtp-username")
		smtpPassword := viper.GetString("smtp-password")
		key := viper.GetString("longy-masterkey")

		smtpHost, smtpPort, err := util.HostAndPort(smtpServer)
		if err != nil {
			return fmt.Errorf("smtp server: %s", err)
		}

		ebSession := eb.CreateSession(authToken, eventID)
		mClient, err := mail.NewClient(smtpHost, smtpPort, smtpUsername, smtpPassword)
		if err != nil {
			return fmt.Errorf("mail client: %s", err)
		}
		mKey, err := mk.NewMasterKey(key)
		if err != nil {
			return fmt.Errorf("master key: %s", err)
		}

		service := rks.NewService(ebSession, mKey, mClient)
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
