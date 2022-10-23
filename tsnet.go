/*
Usage of INET (Uses TSKEY environment variable)
// https://gist.github.com/morgangallant/c77be4c38174f477f8734c09d64638a9
...
if ips := inet.Wait(context.Background()); len(ips) > 0 {
	log.Printf("connected to tailscale (ips=%v)", ips)
} else {
	log.Printf("not connected to tailscale")
}
defer inet.Shutdown()
...
resp, err := inet.HTTPClient().Get("https://tailscale.com")
...
*/

// Package inet is a utility package which creates a Tailscale instance automatically,
// and exposes proper dialers which allow this binary to access a Tailscale network.
package tsproxy

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"time"

	"github.com/pkg/errors"
	"tailscale.com/tsnet"
)

var server = &tsnet.Server{
	AuthKey:   os.Getenv("TSKEY"),
	Logf:      log.Default().Printf,
	Ephemeral: true, // We're only running this in the context of this Go binary.
}

// Status returns the status of the Tailscale server.
func Status(ctx context.Context) (string, []netip.Addr, error) {
	client, err := server.LocalClient()
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get local client")
	}

	resp, err := client.Status(ctx)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get status")
	}

	return resp.BackendState, resp.TailscaleIPs, nil
}

// Wait waits for the server to come online. If the server failed to come online for any reason,
// a nil set of addresses will be returned (indicating that we don't have a Tailscale IP, and are
// thus unable to access machines on the Tailscale network).
func Wait(ctx context.Context) []netip.Addr {
	if _, ok := os.LookupEnv("TSKEY"); !ok {
		return nil
	}
	for {
		status, ips, err := Status(ctx)
		if err != nil {
			return nil
		} else if status == "Running" && len(ips) > 0 {
			return ips
		}
		time.Sleep(time.Second)
	}
}

func Dial() func(ctx context.Context, network, address string) (net.Conn, error) {
	return server.Dial
}

// HTTPClient returns a http.Client which can be used to make outgoing requests to machines on the
// Tailscale network. If the Tailscale server is not running, this will return http.DefaultClient.
func HTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: server.Dial,
		},
	}
}

// Shutdown closes the Tailscale server, if it is running.
func Shutdown() error {
	return server.Close()
}

func noopLogger(format string, args ...any) {}
