package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAttendeeSerialization(t *testing.T) {
	a := NewAttendee("1")

	bz, err := ModuleCdc.MarshalBinaryLengthPrefixed(a)
	require.NoError(t, err)

	var recovered Attendee
	err = ModuleCdc.UnmarshalBinaryLengthPrefixed(bz, &recovered)
	require.NoError(t, err)

	require.Equal(t, a, recovered)
}
