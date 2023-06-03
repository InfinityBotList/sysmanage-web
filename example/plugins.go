package main

import (
	"github.com/infinitybotlist/sysmanage-web/plugins/actions"
	"github.com/infinitybotlist/sysmanage-web/plugins/frontend"
	"github.com/infinitybotlist/sysmanage-web/plugins/nginx"
	"github.com/infinitybotlist/sysmanage-web/plugins/persist"
	"github.com/infinitybotlist/sysmanage-web/plugins/systemd"
	"github.com/infinitybotlist/sysmanage-web/types"
)

var meta = types.ServerMeta{
	Plugins: map[string]types.Plugin{
		"nginx": {
			Init:     nginx.InitPlugin,
			Frontend: "@core",
		},
		"systemd": {
			Init:     systemd.InitPlugin,
			Frontend: "@core",
		},
		// Persist has no frotend, it is a backend plugin
		"persist": {
			Init: persist.InitPlugin,
		},
		"actions": {
			Init:     actions.InitPlugin,
			Frontend: "@core",
		},
		// Frontend has no frontend, it is a backend plugin
		"frontend": {
			Init: frontend.InitPlugin,
		},
	},
	Frontend: types.FrontendConfig{
		FrontendPath:      "frontend", // Use @core here if you want to bootstrap builds from sysmanage itself
		ComponentProvider: "@core",    // Use a custom component provider if you want to modify the components
		CorelibProvider:   "@core",    // Use a custom corelib provider if you want to modify the corelib
		ExtraFiles: map[string]string{
			"frontend/src/lib/images":                 "$lib/images",        // The FrontendPath is irrelevant here, all src's are relative to the root of the project
			"file://frontend/src/lib/images/logo.png": "static/favicon.png", // Since this is specifically a file-file copy, we need to specify this with file://
		},
	},
}
