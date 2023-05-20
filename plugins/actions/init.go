package actions

import "github.com/infinitybotlist/sysmanage-web/types"

func InitPlugin(c *types.PluginConfig) error {
	loadActionsApi(c.Mux)
	return nil
}
