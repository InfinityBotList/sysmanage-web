package server

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/infinitybotlist/sysmanage-web/core/plugins"
	"golang.org/x/exp/slices"
)

var args = map[string]func(){
	"update-frontend": func() {
		// Download sysmanage-web repo
		tmp, err := os.MkdirTemp("", "sysmanage-upd")

		if err != nil {
			panic(err)
		}

		_, err = git.PlainClone(tmp, false, &git.CloneOptions{
			URL: func() string {
				if os.Getenv("GIT_REPO") != "" {
					return os.Getenv("GIT_REPO")
				}

				return "https://github.com/infinitybotlist/sysmanage-web"
			}(),
			Progress: os.Stdout,
			ReferenceName: plumbing.ReferenceName(func() string {
				if os.Getenv("GIT_REF") != "" {
					return os.Getenv("GIT_REF")
				}

				return "refs/heads/main"
			}()),
		})

		if err != nil {
			panic(err)
		}

		// Loop over frontend/src/routes/plugins
		fsd, err := os.ReadDir("frontend/src/routes/plugins")

		if err != nil {
			panic(err)
		}

		for _, f := range fsd {
			if !slices.Contains(plugins.OfficialPlugins, f.Name()) {
				fmt.Println("Skipping custom plugin: " + f.Name())
				continue
			}

			// Remove old plugin
			err = os.RemoveAll("frontend/src/routes/plugins/" + f.Name())

			if err != nil {
				fmt.Println("Failed to remove old plugin: " + f.Name())
				panic(err)
			}

			// Copy new plugin
			err = os.Rename(tmp+"/frontend/src/routes/plugins/"+f.Name(), "frontend/src/routes/plugins/"+f.Name())

			if err != nil {
				fmt.Println("Failed to copy new plugin: " + f.Name())
				panic(err)
			}
		}
	},
}

func parseArgs() {
	arg := os.Args[0]

	if f, ok := args[arg]; ok {
		f()
		return
	}

	fmt.Println("Unknown argument: " + arg)
	os.Exit(1)
}
