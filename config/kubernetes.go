package config

type Kubernetes struct {
	ApiServer string `mapstructure:"api_server" json:"api_server" yaml:"api_server"`
	Timeout   int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	CAFile    string `mapstructure:"ca_file" json:"ca_file" yaml:"ca_file"`
	CertFile  string `mapstructure:"cert_file" json:"cert_file" yaml:"cert_file"`
	KeyFile   string `mapstructure:"key_file" json:"key_file" yaml:"key_file"`
}
