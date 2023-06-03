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
)
