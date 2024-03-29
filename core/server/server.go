package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/infinitybotlist/sysmanage-web/core"
	"github.com/infinitybotlist/sysmanage-web/core/server/cmd"
	"github.com/infinitybotlist/sysmanage-web/core/state"
	"github.com/infinitybotlist/sysmanage-web/types"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/yaml.v3"
)

var frontend fs.FS

var (
	config *types.Config

	// Subbed frontend embed
	serverRootSubbed fs.FS
)

func routeStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api") {
			if state.ServerMeta.FrontendServer != nil {
				fmt.Println("Using external frontend server, query it instead through a simple proxy")
				proxy(w, r)
				return
			}

			serverRoot := http.FS(serverRootSubbed)

			fserve := http.FileServer(serverRoot)
			fserve.ServeHTTP(w, r)
		} else {
			// Serve API
			next.ServeHTTP(w, r)
		}
	})
}

func Init(
	meta types.ServerMeta,
	frontendUi fs.FS,
) {
	core.Init()

	frontend = frontendUi

	state.ServerMeta = meta

	if len(os.Args) > 1 {
		cmd.RunCommand()
		return
	}

	if meta.ConfigVersion < 1 {
		fmt.Println(`
Sysmanage has undergone some big changes between v0 and v1

- Plugin routes are now scoped to the name of the plugin
- Corelib has been updated
- To automatically update corelib, run "sysmanage updatecore"`)
		os.Exit(1)
	}

	// Load config.yaml into Config struct
	file, err := os.Open("config.yaml")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&config)

	if err != nil {
		panic(err)
	}

	state.Config = config

	if meta.FrontendServer != nil {
		fmt.Println("Starting up external frontend server")
		startFrontendServer()
	} else {
		// Create subbed frontend embed
		// Serve frontend
		serverRootSubbed, err = fs.Sub(frontend, "frontend/build")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create wildcard route
	r := chi.NewRouter()

	// A good base middleware stack
	//
	// Load core middleware here
	r.Use(
		middleware.Recoverer,
		middleware.Logger,
		middleware.CleanPath,
		middleware.RealIP,
	)

	// Start loading the plugins
	fmt.Println("Loading plugins...")

	// First run preload scripts
	for _, plugin := range meta.Plugins {
		if plugin.Preload != nil {
			fmt.Println("Running preload action for", plugin.ID)

			err := plugin.Preload(&types.PluginConfig{
				Name:   plugin.ID,
				RawMux: r,
			})

			if err != nil {
				panic(err)
			}
		}
	}

	// Load the other middleware post preload
	r.Use(
		routeStatic,
		middleware.Timeout(30*time.Second),
	)

	for _, plugin := range meta.Plugins {
		fmt.Println("Loading plugin " + plugin.ID)

		if _, ok := state.Config.Plugins[plugin.ID]; !ok {
			panic("Plugin " + plugin.ID + " not found in config.yaml")
		}

		r.Route("/api/"+plugin.ID, func(mr chi.Router) {
			err := plugin.Init(&types.PluginConfig{
				Name:   plugin.ID,
				Mux:    mr,
				RawMux: r,
			})

			if err != nil {
				panic(err)
			}
		})

		state.LoadedPlugins = append(state.LoadedPlugins, plugin.ID)
	}

	if len(state.AuthPlugins) == 0 {
		fmt.Fprintln(os.Stderr, "No auth plugins loaded. For security purposes, please load at least one auth plugin. You can use `authdp` for a reasonably secure auth plugin")
		os.Exit(1)
	}

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
		w.Write([]byte("API endpoint not found..."))
	})

	var port int
	if meta.Port == 0 && config.Port == 0 {
		port = 30010
	}

	if meta.Port != 0 {
		port = meta.Port
	}

	if config.Port != 0 {
		port = config.Port
	}

	// Create server
	s := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}

	// Start server
	fmt.Println("Starting server on port " + strconv.Itoa(port))
	err = s.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
