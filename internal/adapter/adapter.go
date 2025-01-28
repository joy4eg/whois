package adapter

import (
	"context"

	"github.com/cockroachdb/errors"
)

var ErrNotFound = errors.New("adapter not found")

type (
	// Adapter is an interface that defines the behavior for retrieving WHOIS information.
	// It provides a common contract for different WHOIS data sources implementations.
	// Implementations of this interface should handle the specifics of connecting to
	// and querying various WHOIS servers or data sources.
	Adapter interface {
		// Get retrieves WHOIS information for the given host.
		Get(ctx context.Context, host string) (string, error)

		// Server returns the WHOIS server address for the adapter.
		Server() string

		// Name returns the name of the adapter.
		Name() string
	}

	// Options is a map of adapter options.
	Options map[string]string
)

// Create instantiates a new WHOIS adapter based on the specified parameters.
// It returns an implementation of the Adapter interface configured for the requested service.
//
// Parameters:
//   - name: The type of adapter to create ("", "afilias", "arpa", etc)
//   - server: The WHOIS server address to connect to
//   - options: Additional configuration options for the adapter
//
// Returns:
//   - Adapter: The configured WHOIS adapter implementation
//   - error: ErrNotFound if the requested adapter type is not supported
func Create(name string, server string, options Options) (Adapter, error) {
	switch name {
	case "":
		return Standart(server, options), nil
	case "afilias":
		return Afilias(server, options)
	case "arpa":
		return Arpa(server, options)
	case "none":
		return None(server, options)
	case "formatted":
		return Formatted(server, options)
	case "verisign":
		return Verisign(server, options)
	case "web":
		return Web(server, options)
	}
	return nil, errors.Wrap(ErrNotFound, name)
}
