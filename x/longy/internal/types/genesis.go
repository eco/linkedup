package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

// GenesisAttendees is the full array of attendees to initialize
type GenesisAttendees []Attendee

// GenesisScans is the full array of scans to initialize
type GenesisScans []Scan

// GenesisServiceKey is the genesis type for the re-key service
type GenesisServiceKey struct {
	Address       sdk.AccAddress `json:"address"`
	crypto.PubKey `json:"pubkey"`
}

// GenesisPrizes is the full array of prizes for the event
type GenesisPrizes []Prize

// GetID turns the prize tier into its key, assuming tiers are unique
func (p *Prize) GetID() []byte {
	return GetPrizeIDByTier(p.Tier)
}

//GetPrizeIDByTier returns the byte array id for a prize by prefixing its tier
func GetPrizeIDByTier(tier uint) []byte {
	b := []byte(strconv.Itoa(int(tier)))
	return PrizeKey(b)
}

//GetGenesisPrizes returns the array of prizes that we start the chain with
func GetGenesisPrizes() GenesisPrizes {
	return GenesisPrizes{
		Prize{
			Tier:      Tier1,
			RepNeeded: Tier1Rep,
			PrizeText: "Nano Ledger",
			Quantity:  400,
		},
		Prize{
			Tier:      Tier2,
			RepNeeded: Tier2Rep,
			PrizeText: "Key Keeper",
			Quantity:  200,
		},
		Prize{
			Tier:      Tier3,
			RepNeeded: Tier3Rep,
			PrizeText: "Customized SFBW Week Shirt",
			Quantity:  150,
		},
		Prize{
			Tier:      Tier4,
			RepNeeded: Tier4Rep,
			PrizeText: "Customized SFBW Physical Coins",
			Quantity:  100,
		},
		Prize{
			Tier:      Tier5,
			RepNeeded: Tier5Rep,
			PrizeText: "Artwork",
			Quantity:  50,
		},
	}
}
