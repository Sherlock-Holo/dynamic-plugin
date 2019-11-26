package setup

import (
	"github.com/caddyserver/caddy"
	corednsplugin "github.com/coredns/coredns/plugin"
)

type Func func(c *caddy.Controller, next corednsplugin.Handler) (corednsplugin.Handler, error)
