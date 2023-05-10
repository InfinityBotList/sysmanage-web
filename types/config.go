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
