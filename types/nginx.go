package types

type NginxServer struct {
	Names     []string        `yaml:"names" validate:"required"`
	Comment   string          `yaml:"comment" validate:"required"`
	Broken    bool            `yaml:"broken"`
	Locations []NginxLocation `yaml:"locations" validate:"required"`
}

type NginxLocation struct {
	Path  string    `yaml:"path" validate:"required"`
	Proxy string    `yaml:"proxy"`
	Opts  []NginxKV `yaml:"opts"`
}

type NginxKV struct {
	Name  string `yaml:"name" validate:"required"`
	Value string `yaml:"value" validate:"required"`
}

type NginxMeta struct {
	OriginCertPath string `yaml:"origin_cert_path" validate:"required"`
	NginxCertPath  string `yaml:"nginx_cert_path" validate:"required"`
	Common         string `yaml:"common" validate:"required"` // Common config
}

type NginxTemplate struct {
	Servers  []NginxServer
	Meta     NginxMeta
	Domain   string
	CertFile string
	KeyFile  string
}

type NginxYaml struct {
	Servers []NginxServer `yaml:"servers" validate:"required"`
}
