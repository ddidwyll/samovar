package mqtt

import "samovar/lib/cfg"

type config struct {
	cfg.Config
	Host     string `json:"host"`
	Port     int    `json:"port"`
	ClientID string `json:"client_id"`
	Topic    string `json:"topic"`
}

const configName = "mqtt"

func loadConfig() (*config, error) {
	c := &config{}
	c.SetConfigName("mqtt")
	return c, cfg.Load(c)
}
