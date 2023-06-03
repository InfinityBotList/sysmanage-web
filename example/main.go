package main

import (
	"embed"

	"github.com/infinitybotlist/sysmanage-web/core/server"
)

//go:embed all:frontend
var ui embed.FS

func main() {
	// Do something with frontend
	server.Init(meta, ui)
}
