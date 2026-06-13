package client

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/gen"
)

type calc struct {clc.CalcActor}

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{client.calc}")

	c.WatchAs("device_state.t_top", "client_state.t_top")
	c.WatchAs("device_state.t_mid", "client_state.t_mid")
	c.WatchAs("device_state.t_btm", "client_state.t_btm")
	c.WatchAs("device_state.power", "client_state.power")
	c.WatchAs("device_state.power_diff", "client_state.power_diff")
	c.WatchAs("device_state.collect", "client_state.collect")
	c.WatchAs("device_state.press", "client_state.press")
	c.WatchAs("session_state.error", "client_state.error")
	c.WatchAs("session_state.started", "client_state.started")
	c.WatchAs("session_state.device_id", "client_state.device_id")
	c.WatchAs("session_state.devices", "client_state.devices")

	return nil
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	return c.HandleChangeReports(msg)
}
