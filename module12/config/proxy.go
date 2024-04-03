package config

type ProxyConfig struct {
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	Url      string `json:"url" yaml:"url"`
	Port     string `json:"port" yaml:"port"`
	Protocol string `json:"protocol" yaml:"protocol"`
}

func NewProxyConfig() *ProxyConfig {
	return &ProxyConfig{
		Enabled:  false,
		Url:      "",
		Port:     "",
		Protocol: "http",
	}
}
