package config

import (
	"github.com/Sherlock-Holo/dynamic-plugin/config/setup"
	"github.com/caddyserver/caddy"
)

type Config interface {
	Plugins() []string
	GetController(name string) *caddy.Controller
	GetSetupFunc(name string) setup.Func
}
