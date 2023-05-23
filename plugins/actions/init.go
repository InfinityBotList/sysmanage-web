package actions

import (
	"github.com/infinitybotlist/sysmanage-web/plugins/frontend"
	"github.com/infinitybotlist/sysmanage-web/types"
)

func InitPlugin(c *types.PluginConfig) error {
	// Register links
	frontend.AddLink(c, frontend.Link{
		Title:       "Custom Actions",
		Description: "Execute custom actions that are defined by plugins.",
		LinkText:    "Execute An Action",
		Href:        "@root",
	})

	loadActionsApi(c.Mux)
	return nil
}
