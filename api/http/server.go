package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CXeon/traefik-support/api/http/router"
	"github.com/CXeon/traefik-support/internal/modular"
	"github.com/gin-gonic/gin"
)

const (
	ReadTimeout    = 10 * time.Second
	WriteTimeout   = 10 * time.Second
	MaxHeaderBytes = 1 << 20
)

type Server struct {
	srv *http.Server
}

func NewServer(addr string, deps modular.Deps, mods []modular.Module) (*Server, error) {
	engine := gin.New()
	g := router.Default(engine)
	for _, m := range mods {
		if err := m.Build(deps, g); err != nil {
			return nil, fmt.Errorf("module build failed: %w", err)
		}
	}
	return &Server{
		srv: &http.Server{
			Addr:           addr,
			Handler:        engine,
			ReadTimeout:    ReadTimeout,
			WriteTimeout:   WriteTimeout,
			MaxHeaderBytes: MaxHeaderBytes,
		},
	}, nil
}

func (s *Server) Start(_ context.Context) error {
	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
