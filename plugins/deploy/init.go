package deploy

import (
	"errors"
	"sync"

	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/types"
)

const ID = "deploy"

var (
	maxConcurrency = 1
	breakpoint     sync.Mutex
	builds         = map[string]*DeployStatus{}

	deployConfigPath string
)

func InitPlugin(c *types.PluginConfig) error {
	cfgData, err := plugins.GetConfig(c.Name)

	if err != nil {
		return errors.New("Failed to get nginx config: " + err.Error())
	}

	maxConcurrency, err = cfgData.GetInt("max_concurrency")

	if err != nil {
		return err
	}

	deployConfigPath, err = cfgData.GetString("deploy_config_path")

	if err != nil {
		return err
	}

	state.AuthExemptRoutes = append(state.AuthExemptRoutes, "/createDeploy")

	loadDeployApi(c.Mux)

	return nil
}
