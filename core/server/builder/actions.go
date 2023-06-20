package builder

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/types"
)

var info = color.New(color.FgCyan).SprintlnFunc()

//var errorText = color.New(color.FgRed).SprintlnFunc()

type action struct {
	Name string
	Func func()
}

// CopyDir copies the content of src to dst. src should be a full path.
func CopyDir(dst, src string) error {
	return filepath.Walk(src, func(path string, i fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// copy to this path
		outpath := filepath.Join(dst, strings.TrimPrefix(path, src))

		fmt.Print(info("=>", path, "->", outpath))

		if i.IsDir() {
			os.MkdirAll(outpath, i.Mode())
			return nil // means recursive
		}

		// Ensure outpath is also created
		err = os.MkdirAll(filepath.Dir(outpath), 0755)

		if err != nil {
			return err
		}

		// handle irregular files
		if !i.Mode().IsRegular() {
			switch i.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				return os.Symlink(link, outpath)
			}
			return nil
		}

		// copy contents of regular file efficiently

		// open input
		in, _ := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		// create output
		fh, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer fh.Close()

		// make it the same
		fh.Chmod(i.Mode())

		// copy content
		_, err = io.Copy(fh, in)
		return err
	})
}

func CopyProvider(
	sp fs.FS,
	src types.Provider,
	dst string,
) error {
	var srcFs fs.FS

	os.MkdirAll(dst, 0755)

	if src.Provider == "@core" {
		srcFs = sp
	} else {
		srcFs = os.DirFS(src.Provider)
	}

	fs.WalkDir(srcFs, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasPrefix(path, ".svelte-kit") {
			return nil
		}

		for _, ignore := range []string{
			"node_modules",
			".git",
			"build",
		} {
			if strings.Contains(path, ignore) {
				fmt.Print(info("=>", path, "(skipped)"))
				return nil
			}
		}

		fmt.Println("=>", path)
		// Extract each file out
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

		fileInfo, err := f.Stat()

		if err != nil {
			panic(err)
		}

		// handle irregular files
		if !fileInfo.Mode().IsRegular() {
			switch fileInfo.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				return os.Symlink(link, outPath)
			}
			return nil
		}

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

		// Make it the same
		err = nf.Chmod(fileInfo.Mode())

		if err != nil {
			return err
		}

		return nil
	})

	for _, override := range src.Overrides {
		err := CopyProvider(sp, types.Provider{
			Provider: override,
		}, dst)

		if err != nil {
			return err
		}
	}

	return nil
}

var BuildActions = []action{
	{
		Name: "Create build template",
		Func: func() {
			if state.ServerMeta.Frontend.FrontendProvider.Provider == "" {
				panic("Frontend provider is empty")
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

			err = CopyProvider(subbed, state.ServerMeta.Frontend.FrontendProvider, "sm-build")

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
				if plugin.BuildScript != nil {
					fmt.Println("=> Running build script for", name)
					err := plugin.BuildScript(&types.BuildScript{
						RootBuildDir: "sm-build",
						BuildDir:     "sm-build/src/routes/plugins/" + name,
					})

					if err != nil {
						panic(err)
					}
				}

				if plugin.Frontend.Provider == "" {
					continue
				}

				fmt.Println("=> Copying plugin", name)

				subbed, err := fs.Sub(cp, "frontend/coreplugins/"+name)

				if err != nil {
					panic(err)
				}

				dstPath := "src/routes/plugins/" + name
				err = CopyProvider(subbed, plugin.Frontend, "sm-build/"+dstPath)

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

			fmt.Println("=> Copying components to build")

			err = CopyProvider(subbed, state.ServerMeta.Frontend.ComponentProvider, "sm-build/src/lib/components")

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

			fmt.Println("=> Copying corelib to build")

			err = CopyProvider(subbed, state.ServerMeta.Frontend.CorelibProvider, "sm-build/src/lib/corelib")

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

				err := CopyDir(out, src)

				if err != nil {
					panic(err)
				}
			}
		},
	},
	{
		Name: "Build frontend",
		Func: func() {
			cmd := exec.Command("npm", "install")
			cmd.Dir = "sm-build"
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "FORCE_COLOR=true")
			cmd.Env = append(cmd.Env, "COLOR=always")

			err := cmd.Run()

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
			// Create frontend if it doesnt exist
			err := os.MkdirAll("frontend/build", 0755)

			if err != nil {
				panic(err)
			}

			// Remove frontend build dir
			err = os.RemoveAll("frontend/build")

			if err != nil {
				panic(err)
			}

			// Move build dir
			err = os.Rename("sm-build/build", "frontend/build")

			if err != nil {
				panic(err)
			}
		},
	},
}
