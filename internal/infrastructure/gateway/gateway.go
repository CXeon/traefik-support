package gateway

import (
	"context"
	"errors"
	"sync/atomic"

	tilesGateway "github.com/CXeon/tiles/gateway"
	"github.com/CXeon/tiles/gateway/traefik"
)

// Gateway wraps a Traefik gateway client, managing endpoint registration and lifecycle.
type Gateway struct {
	provider *traefik.Provider
	endpoint *tilesGateway.Endpoint
	client   tilesGateway.Client
	running  atomic.Bool
}

func NewGateway(provider *traefik.Provider, endpoint *tilesGateway.Endpoint) *Gateway {
	return &Gateway{
		provider: provider,
		endpoint: endpoint,
	}
}

// Start creates the Traefik client, registers the endpoint, and enables auto-renew.
func (g *Gateway) Start(ctx context.Context) error {
	cli, err := traefik.NewClient(ctx, g.provider)
	if err != nil {
		return err
	}
	if err = cli.Register(ctx, g.endpoint); err != nil {
		_ = cli.Close(ctx)
		return err
	}
	g.client = cli
	g.running.Store(true)
	return nil
}

// Stop deregisters the endpoint and closes the Traefik client connection.
func (g *Gateway) Stop(ctx context.Context) error {
	if !g.running.Swap(false) {
		return nil
	}
	if g.client == nil {
		return nil
	}
	var errs []error
	if err := g.client.Deregister(ctx, g.endpoint); err != nil {
		errs = append(errs, err)
	}
	if err := g.client.Close(ctx); err != nil {
		errs = append(errs, err)
	}
	g.client = nil
	return errors.Join(errs...)
}
