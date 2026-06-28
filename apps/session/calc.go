package session

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/gen"

	"errors"
)

type calc struct{ clc.CalcActor }

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(args ...any) error {
	c.InitCalc("{session.calc}")

	c.Watch(
  	recordCollection,
  	"device_state.collect_value",
  	"device_state.collect_type",
	)

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

	changes := map[string]any{"devices": devices, "ready": "false"}
	err := c.SendRequests("session_state", changes)
	if err != nil {
		return err
	}

	changeDevice := func(args clc.Args, apply clc.ApplyFn) {
		dId := args.MustGet("client_desired_state.device_id").String()
		if device, err := cfg.findById(dId); err != nil {
			apply("session_state.device_id", "")
			apply("session_state.ready", "false")
			apply("session_state.error", "invalid device id")
			apply("session_state.has_error", "true")
		} else {
			apply("session_state.device_id", dId)
			apply("session_state.ready", "true")
			apply("session_state.heat_loss", device.HeatLoss)
			apply("session_state.max_power", device.MaxPower)
			apply("session_state.max_collect", device.MaxCollect)
			apply("session_desired_state.power", device.MaxPower)
		}
	}

	c.Watch(changeDevice, "client_desired_state.device_id")

	return nil
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	return c.HandleChangeReports(msg)
}
