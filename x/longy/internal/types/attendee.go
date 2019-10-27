package types

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	"github.com/tendermint/tendermint/crypto"
)

// Attendee encapsulates attendee information
type Attendee struct {
	ID                 string          `json:"id"`
	Address            sdk.AccAddress  `json:"address"`
	PubKey             crypto.PubKey   `json:"pubKey,omitempty"`
	Name               string          `json:"name,omitempty"`
	UnixTimeSecClaimed int64           `json:"unixTimeSecClaimed,omitempty"` //time when this attendee account was claimed
	Commitment         util.Commitment `json:"commitment,omitempty"`
	Claimed            bool            `json:"claimed,omitempty"`
	Sponsor            bool            `json:"sponsor,omitempty"`
	RsaPublicKey       string          `json:"rsaPublicKey,omitempty"`
	EncryptedInfo      []byte          `json:"encryptedInfo,omitempty"`
	ScanIDs            []string        `json:"scanIds,omitempty"`
	Winnings           []Win           `json:"winnings,omitempty"`
	Rep                uint            `json:"rep,omitempty"`
}

// NewAttendee is the constructor for `Attendee`. New attendees default to 0 rep
// and is unclaimed
func NewAttendee(id string, sponsor bool) Attendee {
	addr := util.IDToAddress(id)

	return Attendee{
		ID:      id,
		Address: addr,
		PubKey:  nil,
		Name:    "",

		Commitment:    nil,
		Claimed:       false,
		EncryptedInfo: []byte{},
		ScanIDs:       []string{},
		Sponsor:       sponsor,
		Rep:           0,
	}
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

//AddWinning adds the winning to the array
func (a *Attendee) AddWinning(winning *Win) (added bool) {
	if winning == nil || winning.Claimed || a.containsWinning(winning) {
		return false
	}
	a.Winnings = append(a.Winnings, *winning)
	return true
}

//ClaimWinning claims the winning by tier and returns true. Returns false if the win is not there
func (a *Attendee) ClaimWinning(tier uint) (claimed bool) {
	for i := range a.Winnings {
		e := &a.Winnings[i]
		if e.Tier == tier {
			if e.Claimed { //already claimed
				return false
			}
			e.Claimed = true
			return true
		}
	}
	return false
}

//GetTier returns the tier group that an attendee is in based on their rep value
func (a *Attendee) GetTier() uint {
	switch {
	case a.Rep < Tier1Rep:
		return Tier0
	case a.Rep < Tier2Rep:
		return Tier1
	case a.Rep < Tier3Rep:
		return Tier2
	case a.Rep < Tier4Rep:
		return Tier3
	case a.Rep < Tier5Rep:
		return Tier4
	default:
		return Tier5
	}
}

func (a *Attendee) containsWinning(won *Win) bool {
	for _, e := range a.Winnings {
		if e.Tier == won.Tier {
			return true
		}
	}
	return false
}

// GetID returns the attendee identifier
//nolint:gocritic
func (a *Attendee) GetID() string {
	return a.ID
}

// GetAddress returns the deterministic address associated with the attendee
//nolint:gocritic
func (a *Attendee) GetAddress() sdk.AccAddress {
	return a.Address
}

// GetRep returns the attendee's current rep
//nolint:gocritic
func (a *Attendee) GetRep() uint {
	return a.Rep
}

// AddRep will add rep to the attendee {
func (a *Attendee) AddRep(r uint) {
	a.Rep += r
}

// CurrentCommitment returns the current commitment associated with this attendee
//nolint:gocritic
func (a *Attendee) CurrentCommitment() util.Commitment {
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

// IsKeyed indicates if this attendee is key'd yet
func (a *Attendee) IsKeyed() bool {
	return a.PubKey != nil
}

// IsClaimed indicates if this attendee is claimed
//nolint:gocritic
func (a *Attendee) IsClaimed() bool {
	return a.Claimed
}

// SetClaimed will mark this attendee as claimed
func (a *Attendee) SetClaimed() {
	a.Claimed = true
}

//Encode encodes a hex byte array
func Encode(src []byte) string {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return fmt.Sprintf("%s", dst)
}

//Decode decodes a string into a hex byte array
func Decode(src string) []byte {

	decoded, err := hex.DecodeString(src)
	if err != nil {
		fmt.Println(err.Error())
	}

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
