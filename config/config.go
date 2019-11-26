package config

import (
	"io"
	"plugin"
	"strings"

	"github.com/Sherlock-Holo/dynamic-plugin/config/setup"
	"github.com/Sherlock-Holo/errors"
	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyfile"
)

var _ Config = new(config)

type config struct {
	key string

	dynamicConfigPath    string
	dynamicConfigContent io.Reader

	names       []string
	paths       map[string]string
	setups      map[string]setup.Func
	tokens      map[string][]caddyfile.Token
	controllers map[string]*caddy.Controller
}

func NewConfig(c *caddy.Controller) (*config, error) {
	cfg, err := parse(c)
	if err != nil {
		return nil, errors.WithMessage(err, "parse config failed")
	}

	if err := cfg.loadDynamicConfig(); err != nil {
		return nil, errors.WithMessage(err, "load dynamic config failed")
	}

	if err := cfg.loadPlugins(); err != nil {
		return nil, errors.WithMessage(err, "load dynamic plugins failed")
	}

	return &cfg, nil
}

func (c *config) GetSetupFunc(name string) setup.Func {
	return c.setups[name]
}

func (c *config) Plugins() []string {
	plugins := make([]string, len(c.names))
	copy(plugins, c.names)

	return plugins
}

func (c *config) GetController(name string) *caddy.Controller {
	tokens, ok := c.tokens[name]
	if !ok {
		return nil
	}

	return &caddy.Controller{
		Key:       c.key,
		Dispenser: caddyfile.NewDispenserTokens(c.dynamicConfigPath, tokens),
	}
}

func (c *config) loadPlugins() error {
	c.setups = make(map[string]setup.Func)

	for _, name := range c.names {
		path := c.paths[name]

		p, err := plugin.Open(path)
		if err != nil {
			return errors.Wrapf(err, "open plugin %s failed", name)
		}

		symbol, err := p.Lookup("Setup")
		if err != nil {
			return errors.Wrapf(err, "lookup `Setup` for plugin %s failed", name)
		}

		f, ok := symbol.(setup.Func)
		if !ok {
			return errors.Wrapf(err, "plugin %s Setup function is not `func(c *caddy.Controller, next corednsplugin.Handler) (corednsplugin.Handler, error)`", name)
		}

		c.setups[name] = f
	}

	return nil
}

func (c *config) loadDynamicConfig() error {
	c.controllers = make(map[string]*caddy.Controller, len(c.names))

	blocks, err := caddyfile.Parse(c.dynamicConfigPath, c.dynamicConfigContent, caddy.ValidDirectives("dns"))
	if err != nil {
		return errors.Wrapf(err, "parse %s failed", c.dynamicConfigPath)
	}

	var block caddyfile.ServerBlock

	for _, b := range blocks {
		if strings.Join(b.Keys, " ") == c.key {
			block = b
			break
		}
	}

	// if not found, block.Keys must be empty
	if len(block.Keys) == 0 {
		return errors.Errorf("key %s not found in %s", c.key, c.dynamicConfigPath)
	}

	for dynamicPluginName, tokens := range block.Tokens {
		if _, ok := c.paths[dynamicPluginName]; !ok {
			// ignore useless tokens
			continue
		}

		ctl := &caddy.Controller{
			Key:       c.key,
			Dispenser: caddyfile.NewDispenserTokens(c.dynamicConfigPath, tokens),
		}

		c.controllers[dynamicPluginName] = ctl
	}

	return nil
}
