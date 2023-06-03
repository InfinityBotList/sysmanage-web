package builder

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/infinitybotlist/sysmanage-web/core/state"
)

var info = color.New(color.FgCyan).SprintlnFunc()
var errorText = color.New(color.FgRed).SprintlnFunc()

type action struct {
	Name string
	Func func()
}

func copy(
	sp fs.FS,
	path string,
	dst string,
) error {
	var srcFs fs.FS

	os.MkdirAll(dst, 0755)

	if strings.HasPrefix(path, "@core") {
		srcFs = sp
	} else {
		srcFs = os.DirFS(path)
	}

	fs.WalkDir(srcFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic("Error while walking " + path + ": " + err.Error())
		}

		var outPath = dst + "/" + path

		if d.IsDir() {
			return os.MkdirAll(outPath, 0755)
		}

		fmt.Print(info("=>", path, "->", outPath))

		// Copy file
		f, err := srcFs.Open(path)

		if err != nil {
			return err
		}

		defer f.Close()

		// Create file
		nf, err := os.Create(outPath)

		if err != nil {
			return err
		}

		_, err = io.Copy(nf, f)

		if err != nil {
			return err
		}

		defer nf.Close()

		fMode, err := f.Stat()

		if err != nil {
			panic(err)
		}

		err = nf.Chmod(fMode.Mode())

		if err != nil {
			panic(err)
		}

		return nil
	})

	if strings.HasPrefix(path, "@core") {
		// Split by :
		parts := strings.Split(path, "::")

		if len(parts) == 1 {
			return nil // No override folder
		} else {
			// Override folder
			overrides := parts[1:]

			for _, override := range overrides {
				err := copy(sp, override, dst)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

var BuildActions = []action{
	{
		Name: "Create build template",
		Func: func() {
			if state.ServerMeta.Frontend.FrontendPath == "" {
				panic("Frontend path is empty")
			}

			cp, ok := state.Assets["frontend"]

			if !ok {
				panic("Core plugins not found")
			}

			subbed, err := fs.Sub(cp, "frontend")

			if err != nil {
				panic(err)
			}

			os.RemoveAll("sm-build")

			err = copy(subbed, state.ServerMeta.Frontend.FrontendPath, "sm-build")

			if err != nil {
				panic(err)
			}

			os.RemoveAll("sm-build/build")
			//os.RemoveAll("sm-build/node_modules") // For reproducibility, remove node_modules and use npm i
		},
	},
	{
		Name: "Copy plugins",
		Func: func() {
			cp, ok := state.Assets["frontend"]

			if !ok {
				panic("Core plugins not found")
			}

			os.MkdirAll("sm-build/src/routes/plugins", 0755)

			for name, plugin := range state.ServerMeta.Plugins {
				if plugin.Frontend == "" {
					continue
				}

				fmt.Println("=> Copying plugin", name)

				subbed, err := fs.Sub(cp, "frontend/coreplugins/"+name)

				if err != nil {
					panic(err)
				}

				dstPath := "src/routes/plugins/" + name
				err = copy(subbed, plugin.Frontend, "sm-build/"+dstPath)

				if err != nil {
					panic(err)
				}
			}
		},
	},
	{
		Name: "Copy components",
		Func: func() {
			cp, ok := state.Assets["frontend"]

			if !ok {
				panic("Components not found")
			}

			subbed, err := fs.Sub(cp, "frontend/src/lib/components")

			if err != nil {
				panic(err)
			}

			err = os.RemoveAll("sm-build/src/lib/components")

			if err != nil {
				panic(err)
			}

			err = os.MkdirAll("sm-build/src/lib/components", 0755)

			if err != nil {
				panic(err)
			}

			if state.ServerMeta.Frontend.ComponentProvider == "@core" && state.ServerMeta.Frontend.FrontendPath != "@core" {
				fmt.Println("=> Removing user-provided components due to @core component provider")

				// Remove user-provided components, replace with core components from compile to ensure consistency
				err = os.RemoveAll(state.ServerMeta.Frontend.FrontendPath + "/src/lib/components")

				if err != nil {
					fmt.Print(errorText("Error while removing user-provided components: " + err.Error()))
				}

				copy(subbed, "@core", state.ServerMeta.Frontend.FrontendPath+"/src/lib/components")
			}

			fmt.Println("=> Copying components to build")

			err = copy(subbed, state.ServerMeta.Frontend.ComponentProvider, "sm-build/src/lib/components")

			if err != nil {
				panic(err)
			}
		},
	},
	{
		Name: "Copy corelib",
		Func: func() {
			cp, ok := state.Assets["frontend"]

			if !ok {
				panic("Corelib not found")
			}

			subbed, err := fs.Sub(cp, "frontend/src/lib/corelib")

			if err != nil {
				panic(err)
			}

			err = os.RemoveAll("sm-build/src/lib/corelib")

			if err != nil {
				panic(err)
			}

			err = os.MkdirAll("sm-build/src/lib/corelib", 0755)

			if err != nil {
				panic(err)
			}

			if state.ServerMeta.Frontend.CorelibProvider == "@core" && state.ServerMeta.Frontend.FrontendPath != "@core" {
				fmt.Println("=> Removing user-provided corelib due to @core corelib provider")

				// Remove user-provided components, replace with core components from compile to ensure consistency
				err = os.RemoveAll(state.ServerMeta.Frontend.FrontendPath + "/src/lib/corelib")

				if err != nil {
					fmt.Print(errorText("Error while removing user-provided corelib: " + err.Error()))
				}

				copy(subbed, "@core", state.ServerMeta.Frontend.FrontendPath+"/src/lib/corelib")
			}

			fmt.Println("=> Copying corelib to build")

			err = copy(subbed, state.ServerMeta.Frontend.CorelibProvider, "sm-build/src/lib/corelib")

			if err != nil {
				panic(err)
			}
		},
	},
	{
		Name: "Copy extra files",
		Func: func() {
			for src, dst := range state.ServerMeta.Frontend.ExtraFiles {
				if src == "" || strings.HasPrefix(src, "@core") {
					panic("Invalid file path, file paths must be absolute: " + src)
				}

				var out = "sm-build/" + dst

				if strings.HasPrefix(dst, "$lib/") {
					out = "sm-build/src/lib/" + strings.Replace(dst, "$lib/", "", 1)
				}

				if strings.HasPrefix(src, "file://") {
					src = strings.Replace(src, "file://", "", 1)

					fmt.Print(info("=>", src, "->", out))

					// Read file
					f, err := os.Open(src)

					if err != nil {
						panic(err)
					}

					// Split file by / to get dir
					parts := strings.Split(out, "/")

					// Get second-last part
					if len(parts) > 1 {
						var dirName string = "sm-build/" + parts[len(parts)-2]

						os.MkdirAll(dirName, 0755)
					}

					// Create file
					cf, err := os.Create(out)

					if err != nil {
						panic(err)
					}

					_, err = io.Copy(cf, f)

					if err != nil {
						panic(err)
					}

					f.Close()
					continue
				}

				copy(nil, src, out)
			}
		},
	},
	{
		Name: "Build frontend",
		Func: func() {
			fmt.Println("=> Hack: patch vite to symlink .bin/vite ../vite/bin/vite.js ")
			err := os.Remove("sm-build/node_modules/.bin/vite")

			if err != nil {
				panic(err)
			}

			err = os.Symlink("../vite/bin/vite.js", "sm-build/node_modules/.bin/vite")

			if err != nil {
				panic(err)
			}

			cmd := exec.Command("npm", "install")
			cmd.Dir = "sm-build"
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "FORCE_COLOR=true")
			cmd.Env = append(cmd.Env, "COLOR=always")

			err = cmd.Run()

			if err != nil {
				panic(err)
			}

			cmd = exec.Command("npm", "run", "build")
			cmd.Dir = "sm-build"
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "FORCE_COLOR=true")
			cmd.Env = append(cmd.Env, "COLOR=always")

			err = cmd.Run()

			if err != nil {
				panic(err)
			}
		},
	},
	{
		Name: "Copy build to out",
		Func: func() {
			err := os.RemoveAll("frontend/build")

			if err != nil {
				panic(err)
			}

			err = os.MkdirAll("frontend/build", 0755)

			if err != nil {
				panic(err)
			}

			err = os.Rename("sm-build/build", "frontend/build")

			if err != nil {
				panic(err)
			}
		},
	},
}
