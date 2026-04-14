package cmd

import (
	"github.com/CXeon/traefik-support/internal/modular"
	forwardauthmodule "github.com/CXeon/traefik-support/internal/modular/forwardauth"
)

func modules() []modular.Module {
	return []modular.Module{
		forwardauthmodule.New(),
	}
}
