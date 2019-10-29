package longy

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

// GenesisState is the genesis struct for the longy module
type GenesisState struct {
	KeyService   GenesisService   `json:"service"`
	BonusService GenesisService   `json:"bonus_service"`
	Attendees    GenesisAttendees `json:"attendees"`
	Scans        GenesisScans     `json:"scans"`
	Prizes       GenesisPrizes    `json:"prizes"`
}

// DefaultGenesisState returns the default genesis struct for the longy module
func DefaultGenesisState() GenesisState {
	return GenesisState{KeyService: GenesisService{}, BonusService: GenesisService{},
		Attendees: GenesisAttendees{}, Scans: GenesisScans{}, Prizes: GenesisPrizes{}}
}

//NewGenesisState returns a genesis object of the state given the input params
func NewGenesisState(service GenesisService, bonusService GenesisService,
	attendees []types.Attendee, scans []types.Scan, prizes types.GenesisPrizes) GenesisState {
	return GenesisState{KeyService: service, BonusService: bonusService,
		Attendees: attendees, Scans: scans, Prizes: prizes}
}

// ValidateGenesis validates that the passed genesis state is valid
//nolint:gocritic
func ValidateGenesis(data GenesisState) error {
	if data.KeyService.Address.Empty() {
		return types.ErrGenesisKeyServiceAddressEmpty("key service address must be set")
	} else if data.BonusService.Address.Empty() {
		return types.ErrGenesisKeyServiceAddressEmpty("bonus address must be set")
	}

	if data.Attendees == nil {
		return types.ErrGenesisAttendeesEmpty("empty genesis attendees")
	}

	if data.Prizes == nil {
		return types.ErrGenesisPrizesEmpty("empty genesis prizes")
	}

	var seenIds = make(map[string]bool)
	for _, a := range data.Attendees {
		if seenIds[a.ID] {
			return fmt.Errorf("duplicate id: %s", a.ID)
		}
		seenIds[a.ID] = true
	}
	return nil
}

// InitGenesis will run module initialization using the genesis state
//nolint:gocritic,gocyclo
func InitGenesis(ctx sdk.Context, k keeper.Keeper, state GenesisState) {
	// create and set of all the attendees and cosmos accounts
	accountKeeper := k.AccountKeeper()
	coinKeeper := k.CoinKeeper()

	// set the attendees
	amount := sdk.NewInt(5000)
	coins := sdk.Coins{sdk.Coin{
		Denom:  ModuleName,
		Amount: amount,
	}}

	// create the master account with coins
	serviceAccount := accountKeeper.GetAccount(ctx, state.KeyService.Address)
	if serviceAccount == nil {
		serviceAccount = accountKeeper.NewAccountWithAddress(ctx, state.KeyService.Address)
		if serviceAccount == nil {
			panic("service account must be set in genesis")
		}
	} else if err := serviceAccount.SetPubKey(state.KeyService.PubKey); err != nil {
		panic(err)
	} else if _, err := coinKeeper.AddCoins(ctx, state.KeyService.Address, coins); err != nil {
		panic(err)
	}
	accountKeeper.SetAccount(ctx, serviceAccount)

	// create the bonus account with coints
	bonusAccount := accountKeeper.GetAccount(ctx, state.BonusService.Address)
	if bonusAccount == nil {
		bonusAccount = accountKeeper.NewAccountWithAddress(ctx, state.BonusService.Address)
		if bonusAccount == nil {
			panic("bonus account must be set in genesis")
		}
	} else if err := bonusAccount.SetPubKey(state.BonusService.PubKey); err != nil {
		panic(err)
	} else if _, err := coinKeeper.AddCoins(ctx, state.BonusService.Address, coins); err != nil {
		panic(err)
	}
	accountKeeper.SetAccount(ctx, bonusAccount)

	// register the bonus and service addresses in the keeper
	if err := k.SetServiceAddress(ctx, state.KeyService.Address); err != nil {
		panic(err)
	}
	if err := k.SetBonusServiceAddress(ctx, state.BonusService.Address); err != nil {
		panic(err)
	}

	for i := range state.Attendees {
		a := &state.Attendees[i]
		if accountKeeper.GetAccount(ctx, a.GetAddress()) == nil {
			account := accountKeeper.NewAccountWithAddress(ctx, a.GetAddress())
			accountKeeper.SetAccount(ctx, account)
			_, e := coinKeeper.AddCoins(ctx, account.GetAddress(), coins)
			if e != nil {
				panic(e)
			}
		}
		k.SetAttendee(ctx, a)
	}

	//set scans
	for i := range state.Scans {
		k.SetScan(ctx, &state.Scans[i])
	}

	//set prizes
	for i := range state.Prizes {
		k.SetPrize(ctx, &state.Prizes[i])
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
//nolint:gocritic
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	service := k.GetService(ctx)
	bonusService := k.GetBonusService(ctx)
	attendees := k.GetAllAttendees(ctx)
	scans := k.GetAllScans(ctx)
	prizes, _ := k.GetPrizes(ctx)
	return NewGenesisState(service, bonusService, attendees, scans, prizes)
}
