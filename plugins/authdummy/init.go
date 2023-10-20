// Simple dummy plugin to act as a dummy.
package authdummy

import (
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/types"
)

const ID = "authdummy"

func InitPlugin(c *types.PluginConfig) error {
	state.AuthPlugins = append(state.AuthPlugins, ID)
	return nil
}
