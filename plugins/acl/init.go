// Package acl allows for access controls on certain API endpoints or certain plugins
package acl

import (
	"github.com/infinitybotlist/sysmanage-web/types"
)

const ID = "acl"

var preloaded bool

func InitPlugin(c *types.PluginConfig) error {
	if !preloaded {
		panic("acl plugin must be preloaded")
	}
	return nil
}

func Preload(c *types.PluginConfig) error {
	c.RawMux.Use(MuxMiddleware)
	preloaded = true
	return nil
}
