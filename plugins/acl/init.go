// Package acl allows for access controls on certain API endpoints or certain plugins
package acl

import (
	"github.com/infinitybotlist/sysmanage-web/types"
)

const ID = "acl"

func InitPlugin(c *types.PluginConfig) error {
	c.RawMux.Use(MuxMiddleware)
	return nil
}
