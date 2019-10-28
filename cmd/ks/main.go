package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	ks "github.com/eco/longy/key-service"
	ksCfg "github.com/eco/longy/key-service/config"
	eb "github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	mk "github.com/eco/longy/key-service/masterkey"
	dbm "github.com/eco/longy/key-service/models"
	"github.com/eco/longy/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	rootCmd.Flags().Int("port", 1337, "port to bind the rekey service")

	rootCmd.Flags().String("longy-chain-id", "longychain", "chain-id of the running longy game")
	rootCmd.Flags().String("longy-restservice", "http://localhost:1317", "scheme://host:port of the full node rest client")
	rootCmd.Flags().String("longy-app-url", "http://localhost:5000", "scheme://host of the client web app")

	// using "master" as the seed
	rootCmd.Flags().String("longy-masterkey",
		"fc613b4dfd6736a7bd268c8a0e74ed0d1c04a959f59dd74ef2874983fd443fca", "hex encoded master private key")

	rootCmd.Flags().String("eventbrite-auth", "", "eventbrite authorization token")
	rootCmd.Flags().Int("eventbrite-event", 0, "id associated with the eventbrite event")

	rootCmd.Flags().String("aws-content-bucket", "linkedup-user-content", "content bucket for user uploads")
	rootCmd.Flags().Bool("email-mock", false, "print email URLs instead of emailing")
	rootCmd.Flags().Bool("localstack", false, "use localstack instead of aws; implies --email-mock")
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

		localstack := viper.GetBool("localstack")

		contentBucket := viper.GetString("aws-content-bucket")

		mockEmail := localstack || viper.GetBool("email-mock")

		longyChainID := viper.GetString("longy-chain-id")
		longyAppURL := viper.GetString("longy-app-url")
		longyRestURL := viper.GetString("longy-restservice")
		ksCfg.SetLongyRestURL(longyRestURL)

		key, err := util.Secp256k1FromHex(viper.GetString("longy-masterkey"))
		if err != nil {
			return fmt.Errorf("masterkey: %s", err)
		}

		/** Eventbrite session **/
		ebSession, err := eb.CreateSession(eventID, authToken)
		if err != nil {
			return err
		}

		awsCfg := session.Must(session.NewSession())

		/** Mail Client **/
		mClient, err := mail.NewMockClient(longyAppURL)
		if !mockEmail {
			mClient, err = mail.NewSESClient(awsCfg, localstack, longyAppURL)
		}
		if err != nil {
			return fmt.Errorf("mail client: %s", err)
		}

		/** Backend DB **/
		db, err := dbm.NewDatabaseContextWithCfg(awsCfg, localstack, contentBucket)
		if err != nil {
			return fmt.Errorf("dynamo: %s", err)
		}

		/** Master key session **/
		mKey, err := mk.NewMasterKey(key, longyChainID)
		if err != nil {
			return fmt.Errorf("master key: %s", err)
		}

		service := ks.NewService(ebSession, &mKey, &db, mClient)
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
