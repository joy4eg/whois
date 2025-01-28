package adapter

import "context"

type verisignAdapter struct {
	server string
}

func (a *verisignAdapter) Get(ctx context.Context, host string) (string, error) {
	return Request(ctx, "="+host, a.server, 0)
}

func (a *verisignAdapter) Server() string {
	return a.server
}

func (*verisignAdapter) Name() string {
	return "verisign"
}

func Verisign(server string, _ Options) (Adapter, error) {
	return &verisignAdapter{
		server: server,
	}, nil
}
