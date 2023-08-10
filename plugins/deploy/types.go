package deploy

import "time"

type DeployMeta struct {
	AllowedIps  []string          `yaml:"allowed_ips"`
	Src         *DeploySource     `yaml:"src"`
	Broken      bool              `yaml:"broken"`
	OutputPath  string            `yaml:"output_path"`
	Script      string            `yaml:"script"`
	Webhooks    []*DeployWebhook  `yaml:"webhooks"`
	Timeout     int               `yaml:"timeout"`
	Env         map[string]string `yaml:"env"`
	ConfigFiles []string          `yaml:"config_files"`
}

type DeploySource struct {
	Type  string `yaml:"type"`
	Url   string `yaml:"url"`
	Token string `yaml:"pat"`
	Ref   string `yaml:"ref"`
}

func (d DeploySource) String() string {
	return d.Type + ": " + d.Url + " (" + d.Ref + ")"
}

type DeployWebhook struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
	Type  string `yaml:"type"`
}

type DeployStatus struct {
	Source    *DeploySource
	CreatedAt time.Time
}

func (d DeployStatus) String() string {
	return d.Source.String() + " - " + d.CreatedAt.Format(time.RFC3339) + " (" + time.Since(d.CreatedAt).String() + ")"
}
