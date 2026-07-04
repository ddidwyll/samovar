package session

import (
	"samovar/lib/cfg"

	"errors"
	"fmt"
)

type device struct {
	Id                 string  `json:"id"`
	HeatLoss           int64   `json:"heat_loss"`
	MaxPower           int64   `json:"max_power"`
	MaxCollect         int64   `json:"max_collect"`
	MinStabilityPeriod int64   `json:"min_stability_period"`
	MinRefluxRatio     float64 `json:"min_reflux_ratio"`
}

type config struct {
	cfg.Config
	Devices []device `json:"devices"`
}

func loadConfig() (*config, error) {
	c := &config{}
	c.SetConfigName("session")
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
	case d.HeatLoss == 0:
		err = "heat_loss is required"
	case d.MaxPower <= 0:
		err = "max_power is required"
	case d.MaxCollect <= 0:
		err = "max_collect is required"
	case d.MinStabilityPeriod <= 0:
		err = "min_stability_period is required"
	case d.MinRefluxRatio <= 0.0:
		err = "min_reflux_ratio is required"
	default:
		return nil
	}

	return errors.New(fmt.Sprintf("Session config error: %s", err))
}

func (c *config) findById(deviceId string) (d device, err error) {
	for _, device := range c.Devices {
		if device.Id == deviceId {
			return device, nil
		}
	}
	return d, errors.New("device not found")
}
