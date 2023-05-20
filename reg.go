package main

import (
	"github.com/infinitybotlist/sysmanage-web/plugins/actions"
	"github.com/infinitybotlist/sysmanage-web/plugins/frontend"
	"github.com/infinitybotlist/sysmanage-web/plugins/nginx"
	"github.com/infinitybotlist/sysmanage-web/plugins/persist"
	"github.com/infinitybotlist/sysmanage-web/plugins/systemd"
	"github.com/infinitybotlist/sysmanage-web/types"
)

type plugin func(c *types.PluginConfig) error

var plugins = map[string]plugin{
	"nginx":    nginx.InitPlugin,
	"systemd":  systemd.InitPlugin,
	"persist":  persist.InitPlugin,
	"actions":  actions.InitPlugin,
	"frontend": frontend.InitPlugin,
}
