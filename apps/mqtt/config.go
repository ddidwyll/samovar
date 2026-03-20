package mqtt

import (
	"samovar/lib/cfg"

	"errors"
	"fmt"
)

type config struct {
	cfg.Config
	Host     string `json:"host"`
	Port     int    `json:"port"`
	ClientID string `json:"client_id"`
	Topic    string `json:"topic"`
}

func loadConfig() (*config, error) {
	c := &config{}
	c.SetConfigName("mqtt")
	cfg.Load(c)
	return c, c.validate()
}

func (c *config) validate() error {
	var key string

	switch {
	case c.Host == "":
		key = "host"
	case c.Port == 0:
		key = "port"
	case c.ClientID == "":
		key = "client_id"
	case c.Topic == "":
		key = "topic"
	default:
		return nil
	}

	err := fmt.Sprintf("Config error: %s is required", key)
	return errors.New(err)
}

func (c *config) url() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *config) clientID() []byte {
	return []byte(c.ClientID)
}

func (c *config) topic() []byte {
	return []byte(c.Topic)
}
