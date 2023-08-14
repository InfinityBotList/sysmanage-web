package deploy

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/plugins/frontend"
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
	// Register links
	frontend.AddLink(c, frontend.Link{
		Title:       "Deploy Management",
		Description: "Manage deployment configs on the system.",
		LinkText:    "Manage Deploys",
		Href:        "@root",
	})

	cfgData, err := plugins.GetConfig(c.Name)

	if err != nil {
		return errors.New("Failed to get deploy config: " + err.Error())
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

	// Also, remove any old stale deploys here too
	go func() {
		if _, err := os.Stat("/tmp/deploys"); err == nil {
			fmt.Println("Removing old deploys")

			err = os.RemoveAll("/tmp/deploys")

			if err != nil {
				panic(err)
			}
		}
	}()

	loadDeployApi(c.Mux)

	return nil
}
