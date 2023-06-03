package types

import "github.com/go-chi/chi/v5"

type Config struct {
	AllowedUsers []string                  `yaml:"allowed_users"`
	DPSecret     string                    `yaml:"dp_secret"`
	DPDisable    bool                      `yaml:"dp_disable"`
	URL          string                    `yaml:"url"`
	GithubPat    string                    `yaml:"github_pat"`
	Plugins      map[string]map[string]any `yaml:"plugins"`
}

type PluginConfig struct {
	Mux  *chi.Mux
	Name string
}

type Plugin struct {
	Init     func(c *PluginConfig) error
	Frontend string // either a git repo, a local dir or @core if its a core plugin
}

type FrontendConfig struct {
	FrontendPath string // the path to the frontend

	ComponentProvider string // either a git repo, a local dir or @core to use the libs from sysmanage itself
	CorelibProvider   string // either a git repo, a local dir or @core to use the libs from sysmanage itself

	// an extra files to load from the corelib provider, key is the path to the file/folder in the src, value is the file/folder to the file in the out
	//
	// Note: the dst is relative to the build folder. If $lib/ is prefix, the dst is relative to the lib folder
	ExtraFiles map[string]string
}

type ServerMeta struct {
	Plugins  map[string]Plugin // List of plugins to load
	Frontend FrontendConfig
}
