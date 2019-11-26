package config

import (
	"bytes"
	"io/ioutil"

	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/plugin"
)

// parse will get dynamic config content and its path,
// and then get all dynamic plugin names and there paths.
func parse(c *caddy.Controller) (cfg config, err error) {
	cfg.paths = make(map[string]string)
	cfg.key = c.Key

	var count int

	for c.Next() {
		count++
		if count > 1 {
			return config{}, plugin.ErrOnce
		}

		for c.NextBlock() {
			switch c.Val() {
			case "conf":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return config{}, c.Err("should has 1 remain argument")
				}

				cfg.dynamicConfigPath = args[0]
				confContent, err := ioutil.ReadFile(cfg.dynamicConfigPath)
				if err != nil {
					return config{}, c.Errf("read conf %s failed: %v", cfg.dynamicConfigPath, err)
				}

				cfg.dynamicConfigContent = bytes.NewReader(confContent)

			default:
				name := c.Val()
				args := c.RemainingArgs()
				if len(args) < 1 {
					return config{}, c.Err("should contain plugin path")
				}

				cfg.names = append(cfg.names, name)
				cfg.paths[name] = args[0]
			}
		}
	}

	return
}
