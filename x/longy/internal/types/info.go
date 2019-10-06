package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Info is the shared info that an attendee has passed to another, The data is encrypted
type Info struct {
	ID   []byte
	Data []byte
}

//NewInfo returns a new initiated info
func NewInfo(sender sdk.AccAddress, receiver sdk.AccAddress, data []byte) (Info, sdk.Error) {
	id, err := GenInfoID(sender, receiver)
	if err != nil {
		return Info{}, err
	}

	if len(data) == 0 {
		return Info{}, ErrDataCannotBeEmpty("data cannot be empty or nil")
	}

	return Info{
		ID:   id,
		Data: data,
	}, nil
}

//GenInfoID generates an id for an info by appending the sender and receiver addresses before prefixing
func GenInfoID(sender sdk.AccAddress, receiver sdk.AccAddress) (id []byte, err sdk.Error) {
	err = CheckSameness(sender, receiver)
	if err != nil {
		return
	}

	val := append(sender, receiver...)
	//append the key so we dont have to do this everywhere
	id = InfoKey(val)
	return
}

//CheckSameness checks the address are not empty or the same
//nolint:interfacer
func CheckSameness(s1 sdk.AccAddress, s2 sdk.AccAddress) (err sdk.Error) {
	if s1.Empty() || s2.Empty() {
		err = ErrAccountAddressEmpty("cannot have empty addresses")
		return
	}

	if s1.Equals(s2) {
		err = ErrAccountsSame("addresses cannot be the same")
		return
	}

	return
}
