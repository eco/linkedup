package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

// GenesisAttendee is the attendee structure in the genesis file
type GenesisAttendee struct {
	ID string `json:"id"`
	//Profile GenesisProfile `json:"profile"`   //gets the full info of the account
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

//GenesisPrize is the genesis type for the prizes
type GenesisPrize struct {
	Tier        int    `json:"tier"`
	ScansNeeded int    `json:"scansNeeded"`
	PrizeText   string `json:"prizeText"`
	Quantity    int    `json:"quantity"`
}

// GenesisPrizes is the full array of prizes for the event
type GenesisPrizes []GenesisPrize

// GetID turns the prize tier into its key, assuming tiers are unique
func (p *GenesisPrize) GetID() []byte {
	return []byte(strconv.Itoa(p.Tier))
}

//GetGenesisPrizes returns the array of prizes that we start the chain with
func GetGenesisPrizes() GenesisPrizes {
	return GenesisPrizes{
		GenesisPrize{
			Tier:        1,
			ScansNeeded: 100,
			PrizeText:   "Nano Ledger",
			Quantity:    400,
		},
		GenesisPrize{
			Tier:        2,
			ScansNeeded: 200,
			PrizeText:   "Key Keeper",
			Quantity:    200,
		},
		GenesisPrize{
			Tier:        3,
			ScansNeeded: 300,
			PrizeText:   "Customized SFBW Week Shirt",
			Quantity:    150,
		},
		GenesisPrize{
			Tier:        4,
			ScansNeeded: 350,
			PrizeText:   "Customized SFBW Physical Coins",
			Quantity:    100,
		},
		GenesisPrize{
			Tier:        5,
			ScansNeeded: 400,
			PrizeText:   "Artwork",
			Quantity:    50,
		},
	}
}
