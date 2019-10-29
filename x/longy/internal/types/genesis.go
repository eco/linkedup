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

// GenesisService is the genesis type for the re-key service
type GenesisService struct {
	Address sdk.AccAddress `json:"address"`
	PubKey  crypto.PubKey  `json:"pubkey"`
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
			Tier:             Tier1,
			RepNeeded:        Tier1Rep,
			PrizeText:        "Pocket Bitcoin White Paper",
			PrizeDescription: "A pocket edition of the Bitcoin white paper by Satoshi Nakamoto",
			Quantity:         Tier1Quantity,
		},
		Prize{
			Tier:      Tier2,
			RepNeeded: Tier2Rep,
			PrizeText: "5 Kong Notes",
			PrizeDescription: "Kong is the first crypto-cash; it combines state-of-the-art secure elements, " +
				"groundbreaking flex PCBs, and time-locked smart contracts. You can keep it in your pocket, or " +
				"spend it just like cash. No transaction fees, no paper trail. Kong does for cash what Bitcoin did " +
				"for money -- it opens the floodgates for the next billion crypto users.",
			Quantity: Tier2Quantity,
		},
		Prize{
			Tier:      Tier3,
			RepNeeded: Tier3Rep,
			PrizeText: "Ledger Nano S",
			PrizeDescription: "Securely hold your crypto assets. The Ledger Nano S is built around the most secure " +
				"type of chip on the market, ensuring optimal security for your crypto.",
			Quantity: Tier3Quantity,
		},
		Prize{
			Tier:      Tier4,
			RepNeeded: Tier4Rep,
			PrizeText: "20 KONG Notes",
			PrizeDescription: "Kong is the first crypto-cash; it combines state-of-the-art secure elements, " +
				"groundbreaking flex PCBs, and time-locked smart contracts. You can keep it in your pocket, or " +
				"spend it just like cash. No transaction fees, no paper trail. Kong does for cash what Bitcoin did " +
				"for money -- it opens the floodgates for the next billion crypto users.",
			Quantity: Tier4Quantity,
		},
		Prize{
			Tier:      Tier5,
			RepNeeded: Tier5Rep,
			PrizeText: "100 KONG Notes",
			PrizeDescription: "Kong is the first crypto-cash; it combines state-of-the-art secure elements, " +
				"groundbreaking flex PCBs, and time-locked smart contracts. You can keep it in your pocket, or " +
				"spend it just like cash. No transaction fees, no paper trail. Kong does for cash what Bitcoin did " +
				"for money -- it opens the floodgates for the next billion crypto users.",
			Quantity: Tier5Quantity,
		}, Prize{
			Tier:             Tier6,
			RepNeeded:        Tier6Rep,
			PrizeText:        "Customized SFBW Physical Coin",
			PrizeDescription: "Limited Edition SFBW 2019 Physical Coin",
			Quantity:         Tier6Quantity,
		}, Prize{
			Tier:      Tier7,
			RepNeeded: Tier7Rep,
			PrizeDescription: "KeyKeep Hardware Wallet\n" +
				"PIN protection against unauthorized use\n" +
				"Additional passphrase protection\n" +
				"Customizable transaction speeds\n" +
				"Limitless wallet addresses on one device",
			Quantity: Tier7Quantity,
		}, Prize{
			Tier:      Tier8,
			RepNeeded: Tier8Rep,
			PrizeText: "Customized SFBW Week Shirt",
			PrizeDescription: "SF Blockchain Week - Assorted Color and Size\n" +
				"Unisex\n" +
				"100% combed and ring-spun cotton (heather colors contain polyester)\n" +
				"Fabric weight: 4.2 oz (142 g/m2)\n" +
				"Shoulder-to-shoulder taping\n" +
				"Side-seamed\n" +
				"Assorted sizes and colors",
			Quantity: Tier8Quantity,
		}, Prize{
			Tier:      Tier9,
			RepNeeded: Tier9Rep,
			PrizeText: "Nakamoto by Cryptograffiti",
			PrizeDescription: "Nakamoto print by cryptograffiti\n" +
				"Signed and numbered, limited-edition of 100\n" +
				"11 in x 13 in (28 cm x 33 cm)\n" +
				"Aqueous pigment fine art, semi-gloss archival rag print",
			Quantity: Tier9Quantity,
		},
	}
}
