package persist

import (
	"errors"
	"fmt"

	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/types"
)

var Author string

var Username string
var Password string

func InitPlugin(c *types.PluginConfig) error {
	cfgData, err := plugins.GetConfig(c.Name)

	if err != nil {
		return errors.New("Failed to get nginx config: " + err.Error())
	}

	Author, err = cfgData.GetString("author")

	if err != nil {
		fmt.Println("INFO: No author set for persist plugin, defaulting to sysmanage-web[auto]")
		Author = "sysmanage-web[auto]"
	}

	Username, err = cfgData.GetString("username")

	if err != nil {
		fmt.Println("INFO: No username set for persist plugin, defaulting to PAT")
		Username = state.Config.GithubPat
	}

	Password, err = cfgData.GetString("password")

	if err != nil {
		fmt.Println("INFO: No password set for persist plugin, defaulting to PAT")
		Password = state.Config.GithubPat
	}

	return nil
}
