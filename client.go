// Package whois provides functionality for querying WHOIS information for domain names.
//
// The package implements a WHOIS client that supports multiple TLD (Top Level Domain) adapters
// and can automatically detect the appropriate WHOIS server based on the domain name.
//
// The client supports loading TLD configuration data from JSON files and provides
// flexible adapter options for different WHOIS server implementations.
//
// Example usage:
//
//	client, err := whois.New()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	result, err := client.Whois(context.Background(), "example.com")
package whois

import (
	"cmp"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/singleflight"

	"github.com/joy4eg/whois/internal/adapter"
	"github.com/joy4eg/whois/internal/data"
)

// client is a whois client that implements the Query interface.
type client struct {
	TLDs  map[string]adapter.Adapter
	SF    singleflight.Group
	Cache struct {
		TTL     time.Duration
		Storage *ristretto.Cache[string, string]
	}
}

// Option is a client option.
type Option func(*client)

func WithCache(ttl time.Duration) Option {
	return func(c *client) {
		cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		})
		if err != nil {
			panic(err)
		}

		c.Cache.TTL = ttl
		c.Cache.Storage = cache
	}
}

// New returns new whois client.
func New(opts ...Option) (Client, error) {
	return newClient(opts...)
}

func newClient(opts ...Option) (*client, error) {
	client := &client{
		TLDs: make(map[string]adapter.Adapter),
	}

	for _, opt := range opts {
		opt(client)
	}

	if err := client.LoadData(); err != nil {
		return nil, errors.Wrap(err, "failed to load WHOIS data")
	}

	return client, nil
}

func (c *client) whois(ctx context.Context, host string, servers ...string) (string, error) {
	v, err, _ := c.SF.Do(host, func() (interface{}, error) {
		if len(servers) == 0 {
			ad, err := c.guess(host)
			if err != nil {
				return "", err
			}
			return ad.Get(ctx, host)
		}

		for _, server := range servers {
			result, err := adapter.Standart(server, nil).Get(ctx, host)
			if err == nil {
				return result, nil
			}
		}
		return "", errors.Errorf("%q: no WHOIS server responded", host)
	})

	if err != nil {
		return "", err
	}

	return v.(string), nil
}

func (c *client) Whois(ctx context.Context, host string, servers ...string) (result string, err error) {
	if c.Cache.Storage != nil {
		if result, ok := c.Cache.Storage.Get(host); ok {
			return result, nil
		}
	}

	result, err = c.whois(ctx, host, servers...)
	if err != nil {
		return "", err
	}

	if c.Cache.Storage != nil {
		c.Cache.Storage.SetWithTTL(host, result, 0, c.Cache.TTL)
	}

	return result, nil
}

// matchesKnownDomain checks if the given host matches any known TLD patterns in the client's TLD map.
// It splits the host by dots and iteratively checks longer TLD patterns to find the most specific match.
//
// For example, given "example.co.uk":
// 1. Checks "example.co.uk"
// 2. Checks "co.uk"
// 3. Checks "uk"
//
// Parameters:
//   - host: The domain name to check against known TLD patterns
//
// Returns:
//   - adapter.Adapter: The matching adapter if found, nil otherwise
func (c *client) matchesKnownDomain(host string) adapter.Adapter {
	parts := strings.Split(host, ".")
	for i := 0; i < len(parts); i++ {
		tld := strings.Join(parts[i:], ".")
		if adapter, ok := c.TLDs[tld]; ok {
			return adapter
		}
	}
	return nil
}

func (c *client) guess(host string) (ad adapter.Adapter, err error) {
	if matchesTLD(host) {
		return adapter.Standart("whois.iana.org", nil), nil
	}

	if ad = c.matchesKnownDomain(host); ad != nil {
		return ad, nil
	}

	return nil, ErrCannotMatchTLD
}

func (c *client) LoadData() error {
	dirs, err := data.Files.ReadDir(".")
	if err != nil {
		return errors.Wrap(err, "failed to read data directory")
	}

	for _, entry := range dirs {
		if entry.IsDir() {
			continue
		}

		switch entry.Name() {
		case "asn16.json", "asn32.json", "ipv4.json", "ipv6.json":
			slog.Debug("unsupported data file", "file", entry.Name())

		case "tld.json":
			slog.Debug("loading TLD data file", "file", entry.Name())
			data, err := data.Files.ReadFile(entry.Name())
			if err != nil {
				return errors.Wrap(err, "failed to read TLD data file")
			}
			if err := c.LoadDataTLD(data); err != nil {
				return errors.Wrap(err, "failed to load TLD data")
			}

		default:
			slog.Debug("skipping unknown data file", "file", entry.Name())
		}
	}
	return nil
}

func (c *client) LoadDataTLD(data []byte) error {
	r := gjson.ParseBytes(data)
	slog.Debug("loading TLD data",
		"version", r.Get("_.schema").String(),
		"updated", r.Get("_.updated").String(),
	)

	for tld, config := range r.Map() {
		if tld == "_" {
			// Skip metadata.
			continue
		}

		var options adapter.Options
		err := json.Unmarshal([]byte(config.Raw), &options)
		if err != nil {
			return errors.Wrapf(err, "%q: failed to unmarshal adapter options", tld)
		}

		ad, err := adapter.Create(
			config.Get("adapter").String(),
			cmp.Or(config.Get("host").String(), config.Get("url").String()),
			options,
		)
		if err != nil {
			return errors.Wrapf(err, "%q: failed to create adapter", tld)
		}
		c.TLDs[tld] = ad
	}
	slog.Debug("TLD data loaded", "count", len(c.TLDs))

	return nil
}

func (c *client) Close() error {
	if c.Cache.Storage != nil {
		c.Cache.Storage.Close()
	}
	return nil
}
