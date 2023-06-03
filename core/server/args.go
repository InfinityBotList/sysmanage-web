package server

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/infinitybotlist/sysmanage-web/core/server/builder"
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
