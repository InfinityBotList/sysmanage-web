package nginx

import (
	"errors"
	"os"

	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/plugins/frontend"
	"github.com/infinitybotlist/sysmanage-web/types"

	"github.com/cloudflare/cloudflare-go"
)

var (
	nginxTemplate    string
	nginxDefinitions string
	cf               *cloudflare.API
)

const ID = "nginx"

func InitPlugin(c *types.PluginConfig) error {
	// Register links
	frontend.AddLink(c, frontend.Link{
		Title:       "Nginx Configuration Management",
		Description: "Manage nginx configurations on the system.",
		LinkText:    "View Nginx Config",
		Href:        "@root",
	})

	frontend.AddLink(c, frontend.Link{
		Title:       "Add Nginx Server",
		Description: "Add a new nginx server block to the system.",
		LinkText:    "Add Server",
		Href:        "@root/new",
	})

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
