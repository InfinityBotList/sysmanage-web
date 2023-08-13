package types

import (
	"github.com/go-chi/chi/v5"
)

type Config struct {
	Plugins map[string]map[string]any `yaml:"plugins"`
	Port    int                       `yaml:"port"`
}

type PluginConfig struct {
	Mux    chi.Router
	RawMux *chi.Mux
	Name   string
}

type Plugin struct {
	ID          string // Note that the id of the plugin should never be changed as it determines API endpoints
	Init        func(c *PluginConfig) error
	BuildScript func(b *BuildScript) error
	Preload     func(c *PluginConfig) error // Function to call on preload. Note that only Name and RawMux are set
	Frontend    Provider
}

type BuildScript struct {
	// Build directory (sm-build/plugins/<plugin name>)
	BuildDir string
	// Root build directory (sm-build)
	RootBuildDir string
}

type Provider struct {
	Provider  string // use @core to use the libs from sysmanage itself, otherwise, specify a local directory
	Overrides []string
}

type FrontendConfig struct {
	FrontendProvider Provider // frontend provider

	ComponentProvider Provider // component provider
	CorelibProvider   Provider // corelib provider

	// an extra files to load from the corelib provider, key is the path to the file/folder in the src, value is the file/folder to the file in the out
	//
	// Note: the dst is relative to the build folder. If $lib/ is prefix, the dst is relative to the lib folder
	ExtraFiles map[string]string
}

type ServerMeta struct {
	ConfigVersion  int
	Port           int
	FrontendServer *FrontendServer // Leave blank to use static frontend
	Plugins        []Plugin        // List of plugins to load
	Frontend       FrontendConfig
}

type FrontendServer struct {
	Host                string
	ExtraHeadersToAllow []string
	Dir                 string   // The directory to serve
	DirAbsolute         bool     // If true, the dir is absolute, otherwise, it's relative to the root of the project
	RunCommand          string   // The command to run to start the server
	ExtraEnv            []string // Extra environment variables to pass to the server, format KEY=value
}
