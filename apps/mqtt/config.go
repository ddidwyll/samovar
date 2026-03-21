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
	var err string

	switch {
	case c.Host == "":
		err = "host is required"
	case c.Port == 0:
		err = "port is required"
	case c.ClientID == "":
		err = "client_id is required"
	case c.Topic == "":
		err = "topic is required"
	default:
		return nil
	}

	return errors.New(fmt.Sprintf("Config error: %s", err))
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
