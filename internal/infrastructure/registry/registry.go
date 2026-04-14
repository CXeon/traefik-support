package registry

import (
	"context"
	"errors"
	"sync/atomic"

	tilesRegistry "github.com/CXeon/tiles/registry"
	"github.com/CXeon/tiles/registry/etcd"
)

// Locator is a narrow interface for service discovery lookups.
// Consumers use this to resolve a service name to a live endpoint.
type Locator interface {
	GetService(ctx context.Context, service string, option ...tilesRegistry.GetServiceOption) (tilesRegistry.Endpoint, error)
}

// WatchConfig holds the service discovery configuration.
type WatchConfig struct {
	// Services is the list of service names to discover and watch.
	// Empty means no discovery is performed.
	Services []string
	// ComProj restricts discovery to specific company+project scopes.
	// nil means use the registered endpoint's own company and project.
	ComProj map[string][]string
}

// Registry wraps an etcd registry client, managing endpoint registration,
// service discovery, and lifecycle.
type Registry struct {
	cfg      etcd.Config
	endpoint *tilesRegistry.Endpoint
	watchCfg WatchConfig
	registry *etcd.Registry
	running  atomic.Bool
}

func NewRegistry(cfg etcd.Config, endpoint *tilesRegistry.Endpoint, watchCfg WatchConfig) *Registry {
	return &Registry{
		cfg:      cfg,
		endpoint: endpoint,
		watchCfg: watchCfg,
	}
}

// Start establishes a connection to etcd, registers the service endpoint,
// then (if watch services are configured) performs an initial Discover
// followed by a Watch to keep the local cache up to date.
func (r *Registry) Start(ctx context.Context) error {
	reg, err := etcd.NewRegistry(r.cfg)
	if err != nil {
		return err
	}
	if err = reg.Register(ctx, r.endpoint); err != nil {
		_ = reg.Close(ctx)
		return err
	}
	r.registry = reg
	r.running.Store(true)

	if len(r.watchCfg.Services) == 0 {
		return nil
	}

	var opts []tilesRegistry.ServiceOption
	if len(r.watchCfg.ComProj) > 0 {
		opts = append(opts, tilesRegistry.WithGetOptComProj(r.watchCfg.ComProj))
	}

	if _, err = reg.Discover(ctx, r.watchCfg.Services, opts...); err != nil {
		_ = r.Stop(ctx)
		return err
	}
	if err = reg.Watch(ctx, r.watchCfg.Services, opts...); err != nil {
		_ = r.Stop(ctx)
		return err
	}

	return nil
}

// GetService implements Locator. It resolves a service name to a live
// endpoint using the local cache populated by Discover and Watch.
func (r *Registry) GetService(ctx context.Context, service string, option ...tilesRegistry.GetServiceOption) (tilesRegistry.Endpoint, error) {
	return r.registry.GetService(ctx, service, option...)
}

// Stop deregisters the service endpoint and closes the etcd connection.
func (r *Registry) Stop(ctx context.Context) error {
	if !r.running.Swap(false) {
		return nil
	}
	if r.registry == nil {
		return nil
	}
	var errs []error
	if err := r.registry.Deregister(ctx, r.endpoint); err != nil {
		errs = append(errs, err)
	}
	if err := r.registry.Close(ctx); err != nil {
		errs = append(errs, err)
	}
	r.registry = nil
	return errors.Join(errs...)
}
