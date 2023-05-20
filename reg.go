package main

import (
	"sysmanage-web/plugins/actions"
	"sysmanage-web/plugins/frontend"
	"sysmanage-web/plugins/nginx"
	"sysmanage-web/plugins/persist"
	"sysmanage-web/plugins/systemd"
	"sysmanage-web/types"
)

type plugin func(c *types.PluginConfig) error

var plugins = map[string]plugin{
	"nginx":    nginx.InitPlugin,
	"systemd":  systemd.InitPlugin,
	"persist":  persist.InitPlugin,
	"actions":  actions.InitPlugin,
	"frontend": frontend.InitPlugin,
}
