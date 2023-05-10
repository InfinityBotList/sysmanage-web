package state

import (
	"context"
	"sync"
	"sysmanage-web/types"

	"github.com/go-playground/validator/v10"
)

var (
	Config *types.Config

	// Mutex to ensure only one large scale operation is running at a time
	LsOp = sync.Mutex{}

	Validator = validator.New()

	LoadedPlugins = []string{}

	Context = context.Background()
)
