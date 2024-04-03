package config

type JaegerConfig struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	EndpointUrl string `json:"endpointUrl" yaml:"endpointUrl"`
}

func NewJaegerConfig() *JaegerConfig {
	return &JaegerConfig{
		Enabled:     false,
		EndpointUrl: "",
	}
}
