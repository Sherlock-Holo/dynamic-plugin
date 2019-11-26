package handler

import (
	"github.com/Sherlock-Holo/dynamic-plugin/config"
	"github.com/Sherlock-Holo/errors"
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	caddy.RegisterPlugin("dynamic-plugin", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(ctl *caddy.Controller) error {
	cfg, err := config.NewConfig(ctl)
	if err != nil {
		return errors.WithMessage(err, "new config failed")
	}

	h, err := newHandler(cfg)
	if err != nil {
		return errors.WithMessage(err, "new handler failed")
	}

	dnsCfg := dnsserver.GetConfig(ctl)

	dnsCfg.AddPlugin(func(next plugin.Handler) plugin.Handler {
		h.nextName = next.Name()
		h.next = next

		return h
	})

	return nil
}
