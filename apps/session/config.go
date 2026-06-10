package session

import (
	"samovar/lib/cfg"

	"errors"
	"fmt"
)

type device struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	HeatLoss   int64  `json:"heat_loss"`
	MaxPower   int64  `json:"max_power"`
	MaxCollect int64  `json:"max_collect"`
}

type config struct {
	cfg.Config
	Devices []device `json:"devices"`
}

func loadConfig() (*config, error) {
	c := &config{}
	c.SetConfigName("client")
	cfg.Load(c)
	return c, c.validate()
}

func (c *config) validate() error {
	if len(c.Devices) == 0 {
		return errors.New("must be at least one device")
	}
	for _, device := range c.Devices {
		if err := device.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (d device) validate() error {
	var err string

	switch {
	case d.Id == "":
		err = "id is required"
	case d.Name == "":
		err = "name is required"
	case d.HeatLoss <= 0:
		err = "heat_loss is required"
	case d.MaxPower <= 0:
		err = "max_power is required"
	case d.MaxCollect <= 0:
		err = "max_collect is required"
	default:
		return nil
	}

	return errors.New(fmt.Sprintf("Session config error: %s", err))
}

func (c *config) findById(deviceId) (device, error) {
	for _, device := range c.Devices {
		if device.Id == deviceId {
			return device, nil
		}
	}
	return nil, errors.New("device not found")
}
