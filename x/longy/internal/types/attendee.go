package types

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
)

// Attendee encapsulates attendee information
type Attendee struct {
	ID      string
	Address sdk.AccAddress

	Commitment util.Commitment
	Claimed    bool
	Sponsor    bool
	EncryptKey string
	ScanIDs    []string

	Rep uint
}

// NewAttendee is the constructor for `Attendee`. New attendees default to 0 rep
// and is unclaimed
func NewAttendee(id string) Attendee {
	addr := util.IDToAddress(id)

	return Attendee{
		ID:      id,
		Address: addr,

		Commitment: nil,
		Claimed:    false,
		ScanIDs:    []string{},

		Rep: 0,
	}
}

// NewAttendeeFromGenesis will create an `Attendee` from `GenesisAttendee`
func NewAttendeeFromGenesis(ga GenesisAttendee) Attendee {
	return NewAttendee(ga.ID)
}

//AddScanID adds the new scan id if it isn't already added
func (a *Attendee) AddScanID(id []byte) (added bool) {
	encoded := Encode(id)
	if len(id) > 0 && !contains(a.ScanIDs, encoded) {
		a.ScanIDs = append(a.ScanIDs, encoded)
		return true
	}
	return false
}

func Encode(src []byte) string {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return fmt.Sprintf("%s", dst)
}

func Decode(src string) []byte {

	decoded, err := hex.DecodeString(src)
	if err != nil {
		//log.Fatal(err)
	}

	//fmt.Printf("%s\n", dst[:n])
	return decoded
}
func contains(s []string, val string) bool {
	for _, a := range s {
		if a == val {
			return true
		}
	}
	return false
}

// GetID returns the attendee identifier
//nolint:gocritic
func (a Attendee) GetID() string {
	return a.ID
}

// GetAddress returns the deterministic address associated with the attendee
//nolint:gocritic
func (a Attendee) GetAddress() sdk.AccAddress {
	return a.Address
}

// GetRep returns the attendee's current rep
//nolint:gocritic
func (a Attendee) GetRep() uint {
	return a.Rep
}

// AddRep will add rep to the attendee {
func (a *Attendee) AddRep(r uint) {
	a.Rep += r
}

// CurrentCommitment returns the current commitment associated with this attendee
//nolint:gocritic
func (a Attendee) CurrentCommitment() util.Commitment {
	return a.Commitment
}

// SetCommitment will set the claim hash for this attendee
//nolint:gocritic
func (a *Attendee) SetCommitment(commitment util.Commitment) {
	bz := make([]byte, commitment.Len())
	copy(bz, commitment.Bytes())
	a.Commitment = bz
}

// ResetCommitment will reset this attendee's commitment to nil
func (a *Attendee) ResetCommitment() {
	a.Commitment = nil
}

// IsClaimed indicates if this attendee is claimed
//nolint:gocritic
func (a Attendee) IsClaimed() bool {
	return a.Claimed
}

// SetClaimed will mark this attendee as claimed
func (a *Attendee) SetClaimed() {
	a.Claimed = true
}
