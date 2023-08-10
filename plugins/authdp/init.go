package authdp

import (
	"errors"

	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/types"
)

const ID = "authdp"

var (
	dpSecret     string
	url          string
	allowedUsers []string
)

var preloaded bool

func InitPlugin(c *types.PluginConfig) error {
	if !preloaded {
		panic("authdp plugin must be preloaded")
	}

	cfgData, err := plugins.GetConfig(c.Name)

	if err != nil {
		return errors.New("Failed to get nginx config: " + err.Error())
	}

	dpSecret, err = cfgData.GetString("dp_secret")

	if err != nil {
		return err
	}

	url, err = cfgData.GetString("url")

	if err != nil {
		return err
	}

	allowedUsers, err = cfgData.GetStringArray("allowed_users")

	if err != nil {
		return err
	}

	state.AuthPlugins = append(state.AuthPlugins, ID)

	return nil
}

func Preload(c *types.PluginConfig) error {
	c.RawMux.Use(DpAuthMiddleware)
	preloaded = true
	return nil
}
