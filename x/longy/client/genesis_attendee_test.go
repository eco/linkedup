package client_test

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	app "github.com/eco/longy"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/client"
	"github.com/eco/longy/x/longy/internal/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMonitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Add Genesis Attendee Test Suite")
}

var _ = Describe("Add Genesis Attendee Tests", func() {
	var cdc *codec.Codec

	BeforeEach(func() {
		cdc = app.MakeCodec()
	})

	It("should fail when service address is invalid account Bech32", func() {
		fakeAddr := "notAnAdreesThatIsBeck"
		appState := getAppState(cdc, fakeAddr, false)
		_, err := client.BuildGenesisState(appState, cdc, []string{fakeAddr})
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.GenesisKeyServiceAccountInvalid))
	})

	It("should fail when service address is not in the accounts genesis", func() {
		realAddr := "cosmos13xkwux707n3j4pcq7249rwf88g7jd8ntzpdz9w"
		appState := getAppState(cdc, realAddr, false)
		_, err := client.BuildGenesisState(appState, cdc, []string{realAddr})
		Expect(err).To(Not(BeNil()))
		Expect(err.Code()).To(Equal(types.GenesisKeyServiceAccountNotPresent))
	})

	It("should succeed when service address is valid", func() {
		realAddr := "cosmos13xkwux707n3j4pcq7249rwf88g7jd8ntzpdz9w"
		appState := getAppState(cdc, realAddr, true)
		_, err := client.BuildGenesisState(appState, cdc, []string{realAddr})
		Expect(err).To(BeNil())
	})

})

func getAppState(cdc *codec.Codec, addrBech string, included bool) map[string]json.RawMessage {
	var accountGenesisState genaccounts.GenesisState
	appState := make(map[string]json.RawMessage)

	if included {
		addr, err := sdk.AccAddressFromBech32(addrBech)
		Expect(err).To(BeNil())
		acc := genaccounts.GenesisAccount{Address: addr}
		Expect(acc).To(Not(BeNil()))
		accountGenesisState = append(accountGenesisState, acc)
	}

	genesisStateBz := cdc.MustMarshalJSON(accountGenesisState)
	appState[genaccounts.ModuleName] = genesisStateBz
	def := longy.DefaultGenesisState()
	appState[longy.ModuleName] = cdc.MustMarshalJSON(def)
	Expect(len(appState)).To(Equal(2))

	return appState
}
