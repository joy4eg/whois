package adapter

import (
	"context"

	"github.com/cockroachdb/errors"
)

type webAdapter struct {
	URL string
}

func (a *webAdapter) Get(_ context.Context, host string) (string, error) {
	return "", errors.Errorf("%q: server does not support WHOIS protocol, try web interface %v", host, a.URL)
}

func (a *webAdapter) Server() string {
	return a.URL
}

func (*webAdapter) Name() string {
	return "web"
}

func Web(url string, _ Options) (Adapter, error) {
	return &webAdapter{
		URL: url,
	}, nil
}
