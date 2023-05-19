package actions

import "sysmanage-web/types"

func InitPlugin(c *types.PluginConfig) error {
	loadActionsApi(c.Mux)
	return nil
}
