package config

type Configuration struct {
	App        App        `mapstructure:"app" json:"app" yaml:"app"`
	Log        Log        `mapstructure:"log" json:"log" yaml:"log"`
	Database   Database   `mapstructure:"database" json:"database" yaml:"database"`
	Kubernetes Kubernetes `mapstructure:"kubernetes" json:"kubernetes" yaml:"kubernetes"`
	Jwt        Jwt        `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}
