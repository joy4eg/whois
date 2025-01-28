package adapter

import (
	"bytes"
	"context"
	"io"
	"net"
	"strconv"

	"github.com/cockroachdb/errors"
)

const DefaultWhoisPort = 43

// Request performs a WHOIS protocol query to a specified host and port.
//
// It establishes a TCP connection to the host:port, sends the query string followed by CRLF,
// and reads the complete response. The connection uses Multipath TCP if available.
//
// Parameters:
//   - ctx: Context for controlling the request lifetime
//   - query: The WHOIS query string to send
//   - host: The WHOIS server hostname or IP address
//   - port: The port number to connect to (defaults to DefaultWhoisPort if 0)
//
// Returns:
//   - string: The complete response from the WHOIS server
//   - error: Any error encountered during the connection, write or read operations
func Request(ctx context.Context, query, host string, port int) (string, error) {
	var d net.Dialer

	if port == 0 {
		port = DefaultWhoisPort
	}

	d.SetMultipathTCP(true)
	conn, err := d.DialContext(ctx, "tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return "", errors.Wrapf(err, "%q: dial failed", host)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(query + "\r\n"))
	if err != nil {
		return "", errors.Wrapf(err, "%q: write failed", host)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		return "", errors.Wrapf(err, "%q: read failed", host)
	}

	return buf.String(), nil
}

type standartAdapter struct {
	server string
}

func (a *standartAdapter) Get(ctx context.Context, host string) (string, error) {
	return Request(ctx, host, a.server, 0)
}

func (a *standartAdapter) Server() string {
	return a.server
}

func (*standartAdapter) Name() string {
	return "standart"
}

func Standart(server string, options Options) Adapter {
	return &standartAdapter{
		server: server,
	}
}
