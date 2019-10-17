package longy

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

// GenesisState is the genesis struct for the longy module
type GenesisState struct {
	KeyService GenesisServiceKey `json:"key_service"`
	Redeem     GenesisRedeemKey  `json:"redeem"`
	Attendees  GenesisAttendees  `json:"attendees"`
	Prizes     GenesisPrizes     `json:"prizes"`
}

// DefaultGenesisState returns the default genesis struct for the longy module
func DefaultGenesisState() GenesisState {
	return GenesisState{KeyService: GenesisServiceKey{}, Redeem: GenesisRedeemKey{}, Attendees: GenesisAttendees{}}
}

// ValidateGenesis validates that the passed genesis state is valid
//nolint:gocritic
func ValidateGenesis(data GenesisState) error {
	if data.KeyService.Address.Empty() {
		return types.ErrGenesisKeyServiceAddressEmpty("key service address must be set")
	}

	if data.Redeem.Address.Empty() {
		return types.ErrGenesisRedeemAddressEmpty("Redeem address must be set")
	}

	if data.Attendees == nil {
		return types.ErrGenesisAttendeesEmpty("empty genesis attendees")
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
//nolint:gocritic
func InitGenesis(ctx sdk.Context, k keeper.Keeper, state GenesisState) {
	// create and set of all the attendees and cosmos accounts
	accountKeeper := k.AccountKeeper()
	coinKeeper := k.CoinKeeper()

	// set the master service account
	masterAccount := accountKeeper.NewAccountWithAddress(ctx, state.KeyService.Address)
	err := masterAccount.SetPubKey(state.KeyService.PubKey) //nolint
	if err != nil {
		panic(err)
	}
	accountKeeper.SetAccount(ctx, masterAccount)
	k.SetMasterAddress(ctx, state.KeyService.Address)

	// set the redeem account
	//redeemAccount := accountKeeper.NewAccountWithAddress(ctx, state.Redeem.Address)
	redeemAccount := accountKeeper.GetAccount(ctx, state.Redeem.Address)
	if redeemAccount == nil {
		panic(fmt.Errorf("the redeem account does not exist"))
	}
	err = k.SetRedeemAccount(ctx, redeemAccount.GetAddress())
	if err != nil {
		panic(err)
	}

	// set the attendees
	amount := sdk.NewInt(5000)
	coins := sdk.Coins{sdk.Coin{
		Denom:  ModuleName,
		Amount: amount,
	}}
	for _, a := range state.Attendees {
		attendee := types.NewAttendeeFromGenesis(a)

		account := accountKeeper.NewAccountWithAddress(ctx, attendee.GetAddress())
		accountKeeper.SetAccount(ctx, account)
		_, e := coinKeeper.AddCoins(ctx, account.GetAddress(), coins)
		if e != nil {
			panic(e)
		}
		//attendee.Address = account.GetAddress()
		k.SetAttendee(ctx, &attendee)
	}

	//set prizes
	for i := range state.Prizes {
		k.SetPrize(ctx, &state.Prizes[i])
	}
}
