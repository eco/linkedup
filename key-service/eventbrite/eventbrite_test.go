package eventbrite

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestProfileRetrievalFromBody(t *testing.T) {
	body := `{
	"foo": "bar",
	"profile": {
		"fooMore": "barMore",
		"fields": "we dont care about",

		"email": "kobe@lakers.com",
		"first_name": "kobe",
		"last_name": "bryant"
		}
	}`
	reader := strings.NewReader(body)

	profile, err := getProfileFromBody(reader)
	expectedProfile := &AttendeeProfile{
		FirstName: "kobe",
		LastName:  "bryant",
		Email:     "kobe@lakers.com",
	}

	require.NoError(t, err)
	require.Equal(t, expectedProfile, profile)
}
