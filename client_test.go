package whois

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientMatchDomains(t *testing.T) {
	t.Parallel()

	client, err := newClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	for _, domain := range []string{"example.com", "example.net", "example.org", "example.lc", "example.gd"} {
		require.NotNil(t, client.matchesKnownDomain(domain))
	}
}
