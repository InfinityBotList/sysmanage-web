package nginx

type NginxServerManage struct {
	Domain string    `validate:"required"`
	Server NginxYaml `validate:"required,dive"`
}

type NginxServer struct {
	ID        string          `yaml:"id" validate:"required"`
	Names     []string        `yaml:"names" validate:"required,min=1"`
	Comment   string          `yaml:"comment" validate:"required"`
	Broken    bool            `yaml:"broken"`
	Locations []NginxLocation `yaml:"locations" validate:"required"`
}

type NginxLocation struct {
	Path  string   `yaml:"path" validate:"required"`
	Proxy string   `yaml:"proxy"`
	Opts  []string `yaml:"opts"`
}

type NginxMeta struct {
	OriginCertPath string `yaml:"origin_cert_path" validate:"required"`
	NginxCertPath  string `yaml:"nginx_cert_path" validate:"required"`
	Common         string `yaml:"common" validate:"required"` // Common config
}

type NginxTemplate struct {
	Servers    []NginxServer
	Meta       NginxMeta
	Domain     string
	CertFile   string
	KeyFile    string
	MetaCommon string
}

type NginxYaml struct {
	Servers []NginxServer `yaml:"servers" validate:"required,dive"`
}

type NginxAPIPublishCert struct {
	Domain string `json:"domain" validate:"required"`
	Cert   string `json:"cert" validate:"required"`
	Key    string `json:"key" validate:"required"`
}
