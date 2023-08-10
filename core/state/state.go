package state

import (
	"context"
	"embed"
	"sync"

	"github.com/infinitybotlist/sysmanage-web/types"

	"github.com/go-playground/validator/v10"
)

var (
	Config *types.Config

	// Plugins
	ServerMeta types.ServerMeta

	// Mutex to ensure only one large scale operation is running at a time
	LsOp = sync.Mutex{}

	Validator = validator.New()

	LoadedPlugins = []string{}

	Context = context.Background()

	Assets map[string]embed.FS

	// Public API. All plugins handling page authentication should add to this array
	AuthPlugins = []string{}

	// Public API. List of routes that should be exempted during authentication. Plugins should add to this array if required
	AuthExemptRoutes = []string{}
)
