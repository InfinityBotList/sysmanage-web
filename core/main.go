package core

import (
	"embed"

	"github.com/infinitybotlist/sysmanage-web/core/state"
)

//go:embed all:frontend
var Frontend embed.FS

func Init() {
	state.Assets = map[string]embed.FS{
		"frontend": Frontend,
	}
}
