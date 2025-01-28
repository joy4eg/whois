package whois

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_extractCreationDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		want    time.Time
		wantErr bool
	}{
		{
			name: "ISO format with Creation Date marker",
			data: []byte("Some text\nCreation Date: 2023-01-02T15:04:05Z\nOther text"),
			want: time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC),
		},
		{
			name: "Simple date with created marker",
			data: []byte("Some text\ncreated: 2023-01-02\nOther text"),
			want: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "RFC3339 with created on marker",
			data: []byte("Some text\ncreated on: 2023-01-02T15:04:05Z\nOther text"),
			want: time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC),
		},
		{
			name: "Timezone offset format",
			data: []byte("Some text\ncreated date: 2023-01-02 15:04:05-0700\nOther text"),
			want: time.Date(2023, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60)),
		},
		{
			name:    "No creation date",
			data:    []byte("Some random text without creation date"),
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "Invalid date format",
			data:    []byte("Creation Date: invalid-date"),
			want:    time.Time{},
			wantErr: false,
		},
		{
			name: "Custom date format",
			data: []byte("Some text\nCreation Date: 2024-02-16T10:30:42.0Z\nOther text"),
			want: time.Date(2024, 2, 16, 10, 30, 42, 0, time.UTC),
		},
		{
			name: "Custom date format",
			data: []byte("Some text\nCreation Date: 2014-03-02T15:04:18Z\nOther text"),
			want: time.Date(2014, 3, 2, 15, 4, 18, 0, time.UTC),
		},
		{
			name: "Custom date format",
			data: []byte("Some text\ncreated:       2014-04-03 07:32:41+03\nOther text"),
			want: time.Date(2014, 4, 3, 7, 32, 41, 0, time.FixedZone("", 3*60*60)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractCreationDate(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
