package config

type GithubConfig struct {
	Token         string       `json:"token" yaml:"token" `
	UseProxy      bool         `json:"use_proxy" yaml:"use_proxy" `
	Proxy         *ProxyConfig `json:"proxy" yaml:"proxy"`
	SSHPrivateKey string       `json:"ssh_private_key" yaml:"ssh_private_key"`
}
