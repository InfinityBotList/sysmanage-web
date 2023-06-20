package server

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/fatih/color"
	"github.com/infinitybotlist/sysmanage-web/core/server/builder"
	"github.com/infinitybotlist/sysmanage-web/core/state"
)

var bold = color.New(color.Bold).SprintFunc()

var args = map[string]func(){
	"build": func() {
		fmt.Println("Sysmanage Frontend Builder")

		for i, a := range builder.BuildActions {
			fmt.Print(bold(fmt.Sprintf("[%d/%d] %s\n", i+1, len(builder.BuildActions), a.Name)))
			a.Func()
		}
	},
	"updatecore": func() {
		files, ok := state.Assets["frontend"]

		if !ok {
			fmt.Println("No frontend assets found")
			os.Exit(1)
		}

		fmt.Println("Updating corelib of frontend")

		err := os.RemoveAll("frontend/src/lib/corelib")

		if err != nil {
			panic(err)
		}

		os.MkdirAll("frontend/src/lib/corelib", 0755)

		fileSubbed, err := fs.Sub(files, "frontend/src/lib/corelib")

		if err != nil {
			panic(err)
		}

		err = builder.CopyProvider(fileSubbed, state.ServerMeta.Frontend.CorelibProvider, "frontend/src/lib/corelib")

		if err != nil {
			panic(err)
		}

		fmt.Println("Done")
	},
}

func parseArgs() {
	arg := os.Args[1]

	if f, ok := args[arg]; ok {
		f()
		return
	}

	fmt.Println("Unknown argument: " + arg)
	os.Exit(1)
}
