package frontend

import (
	"github.com/infinitybotlist/sysmanage-web/types"
)

const ID = "frontend"

func InitPlugin(c *types.PluginConfig) error {
	loadFrontendApi(c.Mux)
	return nil
}
