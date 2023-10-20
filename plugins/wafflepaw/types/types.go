package types

import "time"

// E.g:
// ibl:
//   - PROJ
//
// Represents a project to watch
type Project struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Components  []Component `yaml:"components"`
}

// Represents a component to watch within a project
type Component struct {
	Name         string              `yaml:"name"`
	Integrations []IntegrationConfig `yaml:"integrations"`
	Every        time.Duration       `yaml:"every"`
}

type IntegrationConfig struct {
	// An integration name
	Name string `yaml:"name"`

	// Any integration-specific configuration that should be passed to the integration
	Config map[string]string `yaml:"config"`

	// Variables that should be exported to scope
	//
	// @ means passthrough
	Export []string `yaml:"export"`
}
