package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

// GenesisAttendee is the attendee structure in the genesis file
type GenesisAttendee struct {
	ID      string         `json:"id"`
	Profile GenesisProfile `json:"profile"` //gets the full info of the account
}

// GenesisAttendees is the full array of attendees to initialize
type GenesisAttendees []GenesisAttendee

// GenesisProfile is the profile of the attendee from eventbrite
type GenesisProfile struct {
	Name     string `json:"name"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	JobTitle string `json:"job_title"`
}

// GenesisKeyService is the genesis type for the re-key service
type GenesisKeyService struct {
	Address       sdk.AccAddress `json:"address"`
	crypto.PubKey `json:"pubkey"`
}

// GenesisPrizes is the full array of prizes for the event
type GenesisPrizes []Prize

// GetID turns the prize tier into its key, assuming tiers are unique
func (p *Prize) GetID() []byte {
	return []byte(strconv.Itoa(p.Tier))
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
