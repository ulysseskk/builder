package config

type ProxyConfig struct {
	ProxyUrl string `json:"proxy_url" yaml:"proxy_url" mapstructure:"proxy_url"`
}
