package eventbrite

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestEmailRetrievalFromBody(t *testing.T) {
	body := `{
	"foo": "bar",
	"profile": {
		"fooMore": "barMore",
		"email": "kobe@lakers.com"
		}
	}`
	reader := strings.NewReader(body)

	email, err := getEmailFromBody(reader)
	require.NoError(t, err)
	require.Equal(t, email, "kobe@lakers.com")
}
