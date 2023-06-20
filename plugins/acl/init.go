// Package acl allows for access controls on certain API endpoints or certain plugins
package acl

import (
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/types"
)

func InitPlugin(c *types.PluginConfig) error {
	if state.Config.DPDisable {
		panic("SAFETY VIOLATION: acl plugin requires deployproxy to be enabled")
	}

	pluginLoaded = true

	c.RawMux.Use(MuxMiddleware)
	return nil
}
