// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Hosts  []HostConfig
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}

type HostConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Name string `json:"name,omitempty"`
}

func (h HostConfig) GetMetricsUrl() (url string) {
	url = "http://" + h.Host + ":" + h.Port + "/metrics"
	return
}
