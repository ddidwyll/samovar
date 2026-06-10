package session

import (
	clc "samovar/lib/calc"
	"samovar/lib/val"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type calc struct {
	act.Actor
	config *clc.Config
}

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(args ...any) error {
	c.config = clc.NewConfig(c, "{session.calc}")

	return c.prepareDevices(args...)
}

func (c *calc) prepareDevices(args ...any) error {
	cfg, ok := args[0].(*config)
	if !ok {
		return errors.New("session.calc: invalid device config")
	}

	devices := make([]string, 0, len(cfg.Devices))
	for _, device := range cfg.Devices {
		devices = append(devices, device.Id)
	}

	devicesVal, _ := val.ArrAsArr(devices)
	err := c.config.SendRequest("session_state", "devices", devicesVal)
	if err != nil {
		return err
	}

	calcDevice := func(args clc.Args, apply clc.ApplyFn) {
		dId := args["session_state.device_id"].String()
		if device, err := cfg.findById(dId); err != nil {
			apply("session_state.error", "invalid device id")
		} else {
			apply("session_state.heat_loss", device.HeatLoss)
			apply("session_state.max_power", device.MaxPower)
			apply("session_state.max_collect", device.MaxCollect)
		}
	}

	c.config.Watch(calcDevice, "session_state.device_id")

	return nil
}

func (c *calc) HandleMessage(_ gen.PID, changeReport any) error {
	return c.config.HandleReport(changeReport)
}
