package client

import (
	"samovar/lib/cfg"

	"errors"
	"fmt"
)

type config struct {
	cfg.Config
	Host string `json:"host"`
	Port int64  `json:"port"`
}

func loadConfig() (*config, error) {
	c := &config{}
	c.SetConfigName("client")
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
	default:
		return nil
	}

	return errors.New(fmt.Sprintf("Config error: %s", err))
}
