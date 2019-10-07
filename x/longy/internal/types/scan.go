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
	//D1 is the encrypted data shared by S1 with S2
	D1 []byte
	//D2 is the encrypted data shared by S2 with S1
	D2 []byte
	//Complete is true when both S1 and S2 have posted this scan interaction on-chain
	Complete bool
}

//NewScan creates a new scan and sets its id
func NewScan(s1 sdk.AccAddress, s2 sdk.AccAddress, d1 []byte, d2 []byte) (Scan, sdk.Error) {
	id, err := GenScanID(s1, s2)
	if err != nil {
		return Scan{}, err
	}
	return Scan{
		ID:       id,
		S1:       s1,
		S2:       s2,
		D1:       d1,
		D2:       d2,
		Complete: false,
	}, nil
}

//GenScanID creates the unique id between a scan pair, regardless of the order of the account addresses passed into it
func GenScanID(s1, s2 sdk.AccAddress) (id []byte, err sdk.Error) {
	err = CheckSameness(s1, s2)
	if err != nil {
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
