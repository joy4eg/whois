package whois

import (
	"bytes"
	"strings"
	"time"
)

// Record is a whois record.
type Record struct {
	// Domain name.
	Domain string `json:"domain"`

	// Domain creation date.
	CreatedDate time.Time `json:"created_date"`
}

func extractCreationDate(data []byte) (time.Time, error) {
	markers := []string{"Creation Date:", "created:", "created on:", "created date:", "Domain Registration Date:"}
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05-0700",
		"2006-01-02 15:04:05-07",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
	}

	for _, marker := range markers {
		pos := bytes.Index(data, []byte(marker))
		if pos < 0 {
			continue
		}
		end := bytes.Index(data[pos:], []byte("\n"))
		if end < 0 {
			end = len(data)
		}
		for _, format := range formats {
			t, err := time.Parse(format, strings.TrimSpace(string(data[pos+len(marker):pos+end])))
			if err != nil {
				continue
			}
			return t, nil
		}
	}
	return time.Time{}, nil
}

// ParseRecord parses raw WHOIS data for a given domain and returns a Record structure.
// It processes the raw byte data to extract domain information such as creation date.
//
// Parameters:
//   - domain: The domain name for which the WHOIS data is being parsed
//   - data: Raw WHOIS response data as bytes
//
// Returns:
//   - *Record: A pointer to a Record structure containing the parsed information
//   - error: An error if parsing fails, nil otherwise.
func ParseRecord(domain string, data []byte) (*Record, error) {
	r := new(Record)
	r.Domain = domain

	if t, err := extractCreationDate(data); err == nil {
		r.CreatedDate = t
	}

	return r, nil
}
