package session

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/gen"

	"errors"
)

type calc struct {clc.CalcActor}

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(args ...any) error {
	c.InitCalc("{session.calc}")

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

	changes := map[string]any{"devices": devices, "started": "no"}
	err := c.SendRequests("session_state", changes)
	if err != nil {
		return err
	}

	calcDevice := func(args clc.Args, apply clc.ApplyFn) {
		dId := args["session_state.device_id"].String()
		if device, err := cfg.findById(dId); err != nil {
			apply("session_state.started", "no")
			apply("session_state.error", "invalid device id")
		} else {
			apply("session_state.started", "yes")
			apply("session_state.heat_loss", device.HeatLoss)
			apply("session_state.max_power", device.MaxPower)
			apply("session_state.max_collect", device.MaxCollect)
		}
	}

  changeDeviceId := func(args clc.Args, apply clc.ApplyFn) {
    desiredDeviceId := args["client_desired_state.device_id"]

    if desiredDeviceId.IsNil() {
      return
    }

    apply("session_state.device_id", desiredDeviceId)
  }

	c.Watch(calcDevice, "session_state.device_id")
	c.Watch(changeDeviceId, "client_desired_state.device_id")

	return nil
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	return c.HandleChangeReports(msg)
}
