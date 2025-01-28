package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		query   string
		host    string
		port    int
		wantErr bool
	}{
		{
			name:    "valid request with default port",
			query:   "example.com",
			host:    "whois.iana.org",
			port:    0,
			wantErr: false,
		},
		{
			name:    "invalid host",
			query:   "test",
			host:    "invalid.host",
			port:    43,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			got, err := Request(ctx, tt.query, tt.host, tt.port)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, got)
			t.Logf("response: %s", got)
		})
	}
}
