package nginx

import (
	"errors"
	"os"
	"sysmanage-web/core/plugins"
	"sysmanage-web/types"

	"github.com/cloudflare/cloudflare-go"
)

var (
	nginxTemplate    string
	nginxDefinitions string
	cf               *cloudflare.API
)

func InitPlugin(c *types.PluginConfig) error {
	// Read data/nginxgen/nginx.tmpl
	bytes, err := os.ReadFile("data/nginxgen/nginx.tmpl")

	if err != nil {
		return errors.New("Failed to read nginx.tmpl: " + err.Error())
	}

	nginxTemplate = string(bytes)

	cfgData, err := plugins.GetConfig(c.Name)

	if err != nil {
		return errors.New("Failed to get nginx config: " + err.Error())
	}

	nginxDefinitions, err = cfgData.GetString("nginx_definitions")

	if err != nil {
		return err
	}

	cfApiToken, err := cfgData.GetString("cf_api_token")

	if err == nil {
		api, err := cloudflare.NewWithAPIToken(cfApiToken)

		if err != nil {
			panic(err)
		}

		cf = api

		setupCf()
	}

	loadNginxApi(c.Mux)

	return nil
}
