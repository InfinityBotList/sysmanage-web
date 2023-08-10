package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"github.com/infinitybotlist/sysmanage-web/core/server/cmd/builder"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"golang.org/x/exp/slices"
)

var bold = color.New(color.Bold).SprintFunc()

type Command struct {
	Name        string
	Description string
	Run         func()
}

var Commands []Command

func init() {
	// We use a init function here to avoid circular imports
	Commands = []Command{
		{
			Name:        "build",
			Description: "Build the frontend",
			Run: func() {
				fmt.Println("Sysmanage Frontend Builder")

				for i, a := range builder.BuildActions {
					fmt.Print(bold(fmt.Sprintf("[%d/%d] %s\n", i+1, len(builder.BuildActions), a.Name)))
					a.Func()
				}
			},
		},
		{
			Name:        "updatecore",
			Description: "Update the core frontend library (corelib)",
			Run: func() {
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
		},
		{
			Name:        "updatecomponents",
			Description: "Update the core frontend components (components)",
			Run: func() {
				files, ok := state.Assets["frontend"]

				if !ok {
					fmt.Println("No frontend assets found")
					os.Exit(1)
				}

				fmt.Println("Updating components of frontend")

				err := os.RemoveAll("frontend/src/lib/components")

				if err != nil {
					panic(err)
				}

				os.MkdirAll("frontend/src/lib/components", 0755)

				fileSubbed, err := fs.Sub(files, "frontend/src/lib/components")

				if err != nil {
					panic(err)
				}

				err = builder.CopyProvider(fileSubbed, state.ServerMeta.Frontend.ComponentProvider, "frontend/src/lib/components")

				if err != nil {
					panic(err)
				}

				fmt.Println("Done")
			},
		},
		{
			Name:        "listplugins",
			Description: "List all plugins with their status (official/custom)",
			Run: func() {
				officialPlugins := plugins.GetOfficialPluginList()

				plList := make([]string, 0, len(state.ServerMeta.Plugins))
				for _, pl := range state.ServerMeta.Plugins {
					var plType string = "custom"

					if slices.Contains(officialPlugins, pl.ID) {
						plType = "official"
					}

					plList = append(plList, pl.ID+" ("+plType+")")
				}

				fmt.Println(strings.Join(plList, "\n"))
			},
		},
		{
			Name:        "help",
			Description: "Show this help message",
			Run:         HelpCommand,
		},
	}
}

// Shorthand for “Commands = append(Commands, cmd)“
func AddCommand(cmd Command) {
	Commands = append(Commands, cmd)
}

// Help command
func HelpCommand() {
	fmt.Println("Available commands")
	fmt.Println()

	for _, c := range Commands {
		fmt.Printf("%s - %s\n", c.Name, c.Description)
	}
}

func RunCommand() {
	arg := os.Args[1]

	var found bool

	for _, c := range Commands {
		if c.Name == arg {
			found = true
			c.Run()
			break
		}
	}

	if !found {
		fmt.Println("Unknown argument: " + arg)
		fmt.Println()
		HelpCommand()
	}
}
