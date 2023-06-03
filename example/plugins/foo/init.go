package foo

import (
	"github.com/infinitybotlist/sysmanage-web/plugins/frontend"
	"github.com/infinitybotlist/sysmanage-web/types"
)

func InitPlugin(c *types.PluginConfig) error {
	frontend.AddLink(c, frontend.Link{
		Title:       "Foo",
		Description: "This is a foo link",
		LinkText:    "Foo",
		Href:        "@root",
	})

	return nil
}
