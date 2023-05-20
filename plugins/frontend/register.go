package frontend

import (
	"strings"
	"sysmanage-web/types"
)

var RegisteredLinks []Link

func AddLink(c *types.PluginConfig, link Link) {
	link.Href = strings.ReplaceAll(link.Href, "@root", "/plugins/"+c.Name)
	RegisteredLinks = append(RegisteredLinks, link)
}
