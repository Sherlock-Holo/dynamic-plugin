package handler

import (
	"context"

	"github.com/Sherlock-Holo/dynamic-plugin/config"
	"github.com/Sherlock-Holo/errors"
	corednsplugin "github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type handler struct {
	corednsplugin.Handler

	config.Config

	nextName string
	next     corednsplugin.Handler
}

func newHandler(cfg config.Config) (*handler, error) {
	h := &handler{Config: cfg}

	var next corednsplugin.Handler

	names := h.Plugins()
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]

		setupFunc := h.GetSetupFunc(name)
		ctl := h.GetController(name)

		dynamicPlugin, err := setupFunc(ctl, next)
		if err != nil {
			return nil, errors.Wrapf(err, "setup %s plugin failed", name)
		}

		next = dynamicPlugin
	}

	h.Handler = next

	return h, nil
}

func (h *handler) Name() string {
	return "dynamic-plugin"
}

func (h *handler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (code int, err error) {
	code, err = h.Handler.ServeDNS(ctx, w, r)
	if err != nil {
		return
	}

	if h.next == nil {
		return
	}

	return corednsplugin.NextOrFailure(h.nextName, h.next, ctx, w, r)
}
