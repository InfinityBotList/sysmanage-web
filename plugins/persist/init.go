package persist

import (
	"errors"
	"fmt"

	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/types"
)

var Author string

var Username string
var Password string

var UseTokenAuth bool

const ID = "persist"

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

	UseTokenAuth, err = cfgData.GetBool("use_token_auth")

	if err != nil {
		fmt.Println("INFO: No use_token_auth set for persist plugin, defaulting to false")
	}

	Password, err = cfgData.GetString("password")

	if err != nil {
		return err
	}

	Username, err = cfgData.GetString("username")

	if err != nil {
		if UseTokenAuth {
			fmt.Println("INFO: No username set for persist plugin, defaulting to password")
			Username = Password
		} else {
			return err
		}
	}

	if !UseTokenAuth {
		fmt.Println("INFO: Will use BasicAuth for git operations")
	} else {
		fmt.Println("INFO: Will use TokenAuth for git operations")
	}

	// Always persist to git during initial startup
	go PersistToGit("")

	return nil
}
