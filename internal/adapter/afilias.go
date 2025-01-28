package adapter

import (
	"context"
)

type afiliasAdapter struct {
	server string
}

func (a *afiliasAdapter) Get(ctx context.Context, host string) (string, error) {
	return Request(ctx, host, a.server, 0)
}

func (a *afiliasAdapter) Server() string {
	return a.server
}

func (*afiliasAdapter) Name() string {
	return "afilias"
}

func Afilias(server string, _ Options) (Adapter, error) {
	return &afiliasAdapter{
		server: server,
	}, nil
}
