package main

import (
	"fmt"
	ks "github.com/eco/longy/key-service"
	eb "github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	mk "github.com/eco/longy/key-service/masterkey"
	dbm "github.com/eco/longy/key-service/models"
	"github.com/eco/longy/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"strings"
)

func init() {
	rootCmd.Flags().Int("port", 1337, "port to bind the rekey service")

	rootCmd.Flags().String("longy-chain-id", "longychain", "chain-id of the running longy game")
	rootCmd.Flags().String("longy-restservice", "http://localhost:1317", "scheme://host:port of the full node rest client")
	rootCmd.Flags().String("longy-fullnode", "tcp://localhost:26657", "tcp://host:port the full node for tx submission")

	// using "master" as the seed
	rootCmd.Flags().String("longy-masterkey",
		"fc613b4dfd6736a7bd268c8a0e74ed0d1c04a959f59dd74ef2874983fd443fca", "hex encoded master private key")

	rootCmd.Flags().String("smtp-server", "smtp.gmail.com:587", "host:port of the smtp server")
	rootCmd.Flags().String("smtp-username", "testecolongy@gmail.com", "username of the email account")
	rootCmd.Flags().String("smtp-password", "2019longygame", "password of the email account")

	rootCmd.Flags().String("eventbrite-auth", "", "eventbrite authorization token")
	rootCmd.Flags().Int("eventbrite-event", 0, "id associated with the eventbrite event")

	rootCmd.Flags().String("aws-region", "us-west-2", "aws region for dynamodb")
	rootCmd.Flags().String("aws-dynamo-url", "http://localhost:8000", "dynamodb url")
}

var rootCmd = &cobra.Command{
	Use:          "ks",
	Short:        "key service for the longest chain game",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags()) //nolint
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.AutomaticEnv()

		port := viper.GetInt("port")

		authToken := viper.GetString("eventbrite-auth")
		eventID := viper.GetInt("eventbrite-event")

		awsRegion := viper.GetString("aws-region")
		dynamoURL := viper.GetString("aws-dynamo-url")

		longyChainID := viper.GetString("longy-chain-id")
		longyFullNodeURL := viper.GetString("longy-fullnode")
		longyRestURL := viper.GetString("longy-restservice")

		key, err := util.Secp256k1FromHex(viper.GetString("longy-masterkey"))
		if err != nil {
			return fmt.Errorf("masterkey: %s", err)
		}

		/** Eventbrite session **/
		ebSession := eb.CreateSession(authToken, eventID)
		mClient, err := mail.NewClient(session.Must(session.NewSession(&aws.Config{
			Region: aws.String(awsRegion),
		})))
		if err != nil {
			return fmt.Errorf("mail client: %s", err)
		}

		/** Backend DB **/
		db, err := dbm.NewDatabaseContext(awsRegion, dynamoURL)
		if err != nil {
			return fmt.Errorf("dynamo: %s", err)
		}

		/** Master key session **/
		mKey, err := mk.NewMasterKey(key, longyRestURL, longyFullNodeURL, longyChainID)
		if err != nil {
			return fmt.Errorf("master key: %s", err)
		}

		service := ks.NewService(&ebSession, &mKey, &db, mClient)
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
