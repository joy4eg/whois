package adapter

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
)

type formattedAdapter struct {
	server string
	format string
}

func (a *formattedAdapter) Get(ctx context.Context, host string) (string, error) {
	query := fmt.Sprintf(a.format, host)

	return Request(ctx, query, a.server, 0)
}

func (a *formattedAdapter) Server() string {
	return a.server
}

func (*formattedAdapter) Name() string {
	return "formatted"
}

func Formatted(server string, options Options) (Adapter, error) {
	format, ok := options["format"]
	if !ok || format == "" {
		return nil, errors.Errorf("format option is required")
	}

	return &formattedAdapter{
		server: server,
		format: format,
	}, nil
}
