package adapter

import (
	"context"

	"github.com/cockroachdb/errors"
)

type noneAdapter struct{}

func (*noneAdapter) Get(_ context.Context, host string) (string, error) {
	return "", errors.Errorf("%q: does not have a WHOIS server", host)
}

func (*noneAdapter) Name() string {
	return "none"
}

func (*noneAdapter) Server() string {
	return ""
}

func None(string, Options) (Adapter, error) {
	return &noneAdapter{}, nil
}
