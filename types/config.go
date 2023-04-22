package types

type Config struct {
	AllowedUsers       []string `yaml:"allowed_users"`
	ServiceDefinitions string   `yaml:"service_definitions"` // List of folders where definitions can be found
	DPSecret           string   `yaml:"dp_secret"`
	RedisURL           string   `yaml:"redis_url"`
	DPDisable          bool     `yaml:"dp_disable"`
	ServiceOut         string   `yaml:"service_out"`
	SrvModBypass       []string `yaml:"srvmod_bypass"`
	GithubPat          string   `yaml:"github_pat"`
	URL                string   `yaml:"url"`
}
