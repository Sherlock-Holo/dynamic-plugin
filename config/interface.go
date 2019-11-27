package config

import (
	"github.com/Sherlock-Holo/dynamic-plugin/setup"
	"github.com/caddyserver/caddy"
)

type Config interface {
	Plugins() []string
	GetController(name string) *caddy.Controller
	GetSetupFunc(name string) setup.Func
}
