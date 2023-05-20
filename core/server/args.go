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

		custDir := "tmp-cust"

		err = os.Mkdir(custDir, 0755)

		if err != nil {
			panic(err)
		}

		for _, f := range fsd {
			if !slices.Contains(plugins.OfficialPlugins, f.Name()) {
				// Move out old plugin
				fmt.Println("custom:", f.Name())
				err = os.Rename("frontend/src/routes/plugins/"+f.Name(), custDir+"/"+f.Name())

				if err != nil {
					panic(err)
				}
			}
		}

		// Remove old frontend
		err = os.RemoveAll("frontend")

		if err != nil {
			panic(err)
		}

		// Move in new frontend
		err = os.Rename(tmp+"/frontend", "frontend")

		if err != nil {
			panic(err)
		}

		// Move in custom plugins
		fsd, err = os.ReadDir(custDir)

		if err != nil {
			panic(err)
		}

		for _, f := range fsd {
			err = os.Rename(custDir+"/"+f.Name(), "frontend/src/routes/plugins/"+f.Name())

			if err != nil {
				panic(err)
			}
		}

		err = os.RemoveAll(custDir)

		if err != nil {
			panic(err)
		}

		err = os.RemoveAll(tmp)

		if err != nil {
			panic(err)
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
