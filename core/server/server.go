package server

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
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
	"github.com/infinitybotlist/sysmanage-web/plugins/persist"
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

func ensureDpAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.DPDisable {
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("X-DP-Host") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Host header not found."))
			return
		}

		if r.Header.Get("X-DP-Host") != config.URL {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. Domain rebind detected. Expected " + config.URL + " but got " + r.Header.Get("X-DP-Host")))
			return
		}

		if r.Header.Get("X-DP-UserID") == "" {
			// User is not authenticated
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. Not running under deployproxy?"))
			return
		}

		// Check if user is allowed
		if len(config.AllowedUsers) != 0 {
			var allowed bool

			for _, user := range config.AllowedUsers {
				if user == r.Header.Get("X-DP-UserID") {
					allowed = true
					break
				}
			}

			if !allowed {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized. User not allowed to access this site."))
				return
			}
		}

		// User is possibly allowed
		if r.Header.Get("X-DP-Signature") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Signature header not found."))
			return
		}

		// Check for X-DP-Timestamp
		if r.Header.Get("X-DP-Timestamp") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Timestamp header not found."))
			return
		}

		ts := r.Header.Get("X-DP-Timestamp")

		// Validate DP-Secret next
		h := hmac.New(sha512.New, []byte(config.DPSecret))
		h.Write([]byte(ts))
		h.Write([]byte(r.Header.Get("X-DP-UserID")))
		hexed := hex.EncodeToString(h.Sum(nil))

		if r.Header.Get("X-DP-Signature") != hexed {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. Signature from deployproxy mismatch"))
			return
		}

		// Check if timestamp is valid
		timestamp, err := strconv.ParseInt(ts, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Timestamp is not a valid integer."))
			return
		}

		if time.Now().Unix()-timestamp > 10 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized. X-DP-Timestamp is too old."))
			return
		}

		// User is allowed
		next.ServeHTTP(w, r)
	})
}

func routeStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api") {
			if state.ServerMeta.FrontendServer != nil {
				fmt.Println("Using external frontend server, query it instead through a simple proxy")
				proxy(w, r)
				return
			}

			serverRoot := http.FS(serverRootSubbed)

			// Get file extension
			fileExt := ""
			if strings.Contains(r.URL.Path, ".") {
				fileExt = r.URL.Path[strings.LastIndex(r.URL.Path, "."):]
			}

			if fileExt == "" && r.URL.Path != "/" {
				r.URL.Path += ".html"
			}

			if r.URL.Path != "/" {
				r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			}

			var checkPath = r.URL.Path

			if r.URL.Path == "/" {
				checkPath = "/index.html"
			}

			// Check if file exists
			f, err := serverRoot.Open(checkPath)

			if err != nil {
				w.Header().Set("Location", "/404?from="+r.URL.Path)
				w.WriteHeader(http.StatusFound)
				return
			}

			f.Close()

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
	r.Use(
		middleware.Recoverer,
		middleware.Logger,
		middleware.CleanPath,
		middleware.RealIP,
		ensureDpAuth,
		routeStatic,
		middleware.Timeout(30*time.Second),
	)

	// Start loading the plugins
	fmt.Println("Loading plugins...")

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

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
		w.Write([]byte("API endpoint not found..."))
	})

	if meta.Port == 0 {
		meta.Port = 30010
	}

	// Create server
	s := &http.Server{
		Addr:    ":" + strconv.Itoa(meta.Port),
		Handler: r,
	}

	// Always persist to git during initial startup
	go persist.PersistToGit("")

	// Also, remove any old stale deploys here too
	go func() {
		if _, err := os.Stat("deploys"); err == nil {
			fmt.Println("Removing old deploys")

			err = os.RemoveAll("deploys")

			if err != nil {
				panic(err)
			}
		}
	}()

	// Start server
	fmt.Println("Starting server on port " + strconv.Itoa(meta.Port))
	err = s.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
