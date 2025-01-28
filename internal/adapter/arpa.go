package adapter

import (
	"context"

	"github.com/cockroachdb/errors"
)

type arpaAdapter struct{}

func (*arpaAdapter) Get(ctx context.Context, host string) (string, error) {
	return "", errors.Errorf("%q: not implemented", host)
}

func (a *arpaAdapter) Server() string {
	return ""
}

func (*arpaAdapter) Name() string {
	return "arpa"
}

func Arpa(string, Options) (Adapter, error) {
	return &arpaAdapter{}, nil
}
