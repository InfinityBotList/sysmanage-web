package systemd

type ServiceManage struct {
	Service TemplateYaml
	Status  string
	ID      string
}

// Struct used to create a service
type CreateTemplate struct {
	Name    string       `yaml:"name" validate:"required"`
	Service TemplateYaml `yaml:"service" validate:"required"`
}

// Struct used to delete a service
type DeleteTemplate struct {
	Name string `yaml:"name" validate:"required"`
}

// Defines a template which is any FILENAME.yaml where FILENAME != _meta
type TemplateYaml struct {
	Command     string `yaml:"cmd" validate:"required"`         // ExecStart in systemd
	Directory   string `yaml:"dir" validate:"required"`         // WorkingDirectory in systemd
	Target      string `yaml:"target" validate:"required"`      // PartOf in systemd
	Description string `yaml:"description" validate:"required"` // Description in systemd
	After       string `yaml:"after" validate:"required"`       // After in systemd
	Broken      bool   `yaml:"broken"`                          // Does the service even work?
	User        string `yaml:"user"`                            // User in systemd, defaults to root if unset
	Group       string `yaml:"group"`                           // Group in systemd, if-else it defaults to root

	// Only used by sysmanage
	Git *Git `yaml:"git,omitempty" json:"Git,omitempty"`
}

// Defines a git integration
type Git struct {
	Repo          string            `yaml:"repo" validate:"required"`           // Git repo to clone
	Ref           string            `yaml:"ref" validate:"required"`            // e.g. refs/heads/priv-serverlist
	BuildCommands []string          `yaml:"build_commands" validate:"required"` // Commands to run after cloning
	Env           map[string]string `yaml:"env" validate:"required"`            // Environment variables to set
	AllowDirty    bool              `yaml:"allow_dirty"`                        // Allow dirty builds
	ConfigFiles   []string          `yaml:"config_files"`                       // Config files to copy from repo to service, only needed if AllowDirty is false
}

// Defines metadata which is _meta.yaml
type MetaYAML struct {
	Targets []MetaTarget `yaml:"targets" validate:"required"` // List of targets to generate
}

// Defines a target in _meta.yaml:targets
type MetaTarget struct {
	Name        string `yaml:"name" validate:"required"`        // Name of target file
	Description string `yaml:"description" validate:"required"` // Directory to place target file
}
