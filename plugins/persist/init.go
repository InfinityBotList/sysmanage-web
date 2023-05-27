package persist

import (
	"errors"
	"fmt"

	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/types"
)

var Username string

func InitPlugin(c *types.PluginConfig) error {
	cfgData, err := plugins.GetConfig(c.Name)

	if err != nil {
		return errors.New("Failed to get nginx config: " + err.Error())
	}

	Username, err = cfgData.GetString("username")

	if err != nil {
		fmt.Println("WARNING: No username set for persist plugin, defaulting to sysmanage-web[auto]")
		Username = "sysmanage-web[auto]"
	}

	return nil
}
