package systemd

type ServiceManage struct {
	Service    *TemplateYaml
	RawService *RawService // Only set when service is not the typical yaml file format
	Status     string
	ID         string
}

type RawService struct {
	Body     string
	FileName string
}

// Struct used to create a service
type CreateTemplate struct {
	Name       string        `yaml:"name"`
	Service    *TemplateYaml `yaml:"service"`
	RawService *RawService   `yaml:"raw_service"` // Only set when service is not the typical yaml file format
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
