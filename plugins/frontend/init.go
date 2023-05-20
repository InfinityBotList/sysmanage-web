package frontend

import "sysmanage-web/types"

func InitPlugin(c *types.PluginConfig) error {
	loadFrontendApi(c.Mux)
	return nil
}
