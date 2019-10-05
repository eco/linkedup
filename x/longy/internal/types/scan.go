package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Scan represents an unique scan between two parties, can only be one scan between the same parties
type Scan struct {
	//ID is the key we use to store this struct in the keyStore, it is unique per S1-S2 pair
	ID []byte
	//S1 is the scanner that initiates the scan
	S1 sdk.AccAddress
	//S2 is the person who's QR code is scanned
	S2 sdk.AccAddress
	//Complete is true when both S1 and S2 have posted this scan interaction on-chain
	Complete bool
}

//NewScan creates a new scan and sets its id
func NewScan(s1 sdk.AccAddress, s2 sdk.AccAddress) (*Scan, sdk.Error) {
	id, err := GenID(s1, s2)
	if err != nil {
		return nil, err
	}
	return &Scan{
		ID:       id,
		S1:       s1,
		S2:       s2,
		Complete: false,
	}, nil
}

//GenID creates the unique id between a scan pair, regardless of the order of the account addresses passed into it
func GenID(s1, s2 sdk.AccAddress) (id []byte, err sdk.Error) {
	if s1.Empty() || s2.Empty() {
		err = ErrAccountAddressEmpty("cannot create a scan where an address is empty")
		return
	}

	if s1.Equals(s2) {
		err = ErrScanAccountsSame("cannot create a scan where addresses are the same")
		return
	}

	val := bytes.Compare(s1, s2)

	//nolint:gocritic
	if val > 0 {
		id = append(s1, s2...)
	} else {
		id = append(s2, s1...)
	}
	//append the key so we dont have to do this everywhere
	id = ScanKey(id)
	return
}
