package whois

import (
	"context"
	"errors"
	"io"
)

var (
	ErrCannotMatchTLD = errors.New("cannot match TLD")
)

// Client is a whois client.
type Client interface {
	// Whois returns the result of a whois query for the given host.
	// The servers parameter is a list of whois servers to try in order, or nil to use the default list.
	// The result is the raw whois output.
	Whois(ctx context.Context, host string, servers ...string) (result string, err error)

	io.Closer
}
