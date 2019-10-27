package longy_test

import (
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/eco/longy/x/longy/internal/types"
	"github.com/eco/longy/x/longy/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	crypto "github.com/tendermint/tendermint/crypto/secp256k1"
)

var _ = Describe("Genesis Tests", func() {
	BeforeEach(func() {
		BeforeTestRun()
	})

	Context("ValidateGenesis", func() {
		It("should fail to validate when service is not set", func() {
			state := longy.GenesisState{
				KeyService: longy.GenesisService{},
				Attendees:  nil,
				Scans:      nil,
				Prizes:     nil,
			}
			err := longy.ValidateGenesis(state)
			Expect(err).ToNot(BeNil())
		})

		It("should fail to validate when attendees are not set", func() {
			state := longy.GenesisState{
				KeyService: longy.GenesisService{
					Address: util.IDToAddress("asdf"),
					PubKey:  crypto.GenPrivKeySecp256k1([]byte("service")).PubKey(),
				},
				Attendees: nil,
				Scans:     nil,
				Prizes:    nil,
			}
			err := longy.ValidateGenesis(state)
			Expect(err).ToNot(BeNil())
		})

		It("should fail to validate when prizes are not set", func() {
			state := longy.GenesisState{
				KeyService: longy.GenesisService{
					Address: util.IDToAddress("asdf"),
					PubKey:  crypto.GenPrivKeySecp256k1([]byte("service")).PubKey(),
				},
				Attendees: longy.GenesisAttendees{},
				Scans:     nil,
				Prizes:    nil,
			}
			err := longy.ValidateGenesis(state)
			Expect(err).ToNot(BeNil())
		})

		It("should succeed to validate when requirements are set", func() {
			state := longy.GenesisState{
				KeyService: longy.GenesisService{
					Address: util.IDToAddress("asdf"),
					PubKey:  crypto.GenPrivKeySecp256k1([]byte("service")).PubKey(),
				},
				Attendees: longy.GenesisAttendees{},
				Scans:     nil,
				Prizes:    types.GetGenesisPrizes(),
			}
			err := longy.ValidateGenesis(state)
			Expect(err).To(BeNil())
		})
	})

	Context("ExportGenesis", func() {
		const (
			q1 = "asdf"
			q2 = "1234"
		)
		It("should export the key service", func() {
			//Create service
			service := longy.GenesisService{
				Address: util.IDToAddress("asdf"),
				PubKey:  crypto.GenPrivKeySecp256k1([]byte("service")).PubKey(),
			}
			account := keeper.AccountKeeper().NewAccountWithAddress(ctx, service.Address)
			keeper.AccountKeeper().SetAccount(ctx, account)

			err := keeper.SetServiceAddress(ctx, service.Address)
			Expect(err).To(BeNil())
			acc := keeper.AccountKeeper().GetAccount(ctx, service.Address)
			e := acc.SetPubKey(service.PubKey)
			Expect(e).To(BeNil())
			keeper.AccountKeeper().SetAccount(ctx, acc)

			genesis := longy.ExportGenesis(ctx, keeper)
			Expect(genesis.KeyService).ToNot(BeNil())
			Expect(genesis.KeyService.PubKey).To(Equal(service.PubKey))
			Expect(genesis.KeyService.Address.Equals(service.Address)).To(BeTrue())
		})

		It("should export the attendees", func() {
			utils.AddAttendeeToKeeper(ctx, &keeper, q1, true, false)
			utils.AddAttendeeToKeeper(ctx, &keeper, q2, true, false)

			genesis := longy.ExportGenesis(ctx, keeper)
			Expect(len(genesis.Attendees)).To(Equal(2))
		})

		It("should export the scans", func() {
			s1 := util.IDToAddress("1234")
			s2 := util.IDToAddress("asdf")
			d1 := []byte{1}
			d2 := []byte{2}
			scan, err := types.NewScan(s1, s2, d1, d2, 1, 2)
			Expect(err).To(BeNil())

			scan1, err := types.NewScan(util.IDToAddress("2222"), util.IDToAddress("33333"),
				[]byte{3}, []byte{4}, 5, 6)
			Expect(err).To(BeNil())
			keeper.SetScan(ctx, scan)
			keeper.SetScan(ctx, scan1)

			genesis := longy.ExportGenesis(ctx, keeper)
			Expect(len(genesis.Scans)).To(Equal(2))
		})

		It("should export the prizes", func() {
			prizes := types.GetGenesisPrizes()
			for i := range prizes {
				keeper.SetPrize(ctx, &prizes[i])
			}
			genesis := longy.ExportGenesis(ctx, keeper)
			Expect(len(genesis.Prizes)).To(Equal(len(prizes)))
			Expect(genesis.Prizes[0].Quantity).To(Equal(prizes[0].Quantity))
		})

	})
})
