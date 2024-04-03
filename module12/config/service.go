package config

type ServiceConfig struct {
	Name string `json:"name" yaml:"name"`
}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		Name: "",
	}
}
