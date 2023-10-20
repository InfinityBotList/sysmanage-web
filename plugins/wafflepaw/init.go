package wafflepaw

import (
	"encoding/json"
	"errors"
	"os"

	wtypes "github.com/infinitybotlist/sysmanage-web/plugins/wafflepaw/types"
	"github.com/infinitybotlist/sysmanage-web/types"
)

var (
	config  []*wtypes.Project
	secrets map[string]map[string]string
)

const ID = "wafflepaw"

func InitPlugin(c *types.PluginConfig) error {
	bytes, err := os.ReadFile("data/wafflepaw/projects.yaml")

	if err != nil {
		return errors.New("Failed to read projects.yaml: " + err.Error())
	}

	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return errors.New("Failed to unmarshal projects.yaml: " + err.Error())
	}

	bytes, err = os.ReadFile("data/wafflepaw/secrets.yaml")

	if err != nil {
		return errors.New("Failed to read secrets.yaml: " + err.Error())
	}

	err = json.Unmarshal(bytes, &secrets)

	if err != nil {
		return errors.New("Failed to unmarshal secrets.yaml: " + err.Error())
	}

	startWafflepawCluster()

	return nil
}
