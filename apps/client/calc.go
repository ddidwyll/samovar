package client

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type calc struct {
	act.Actor
	config *clc.Config
}

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.config = clc.NewConfig(c, "{client.calc}")

	c.config.WatchAs("device_state.t_top", "client_state.t_top")
	c.config.WatchAs("device_state.t_mid", "client_state.t_mid")
	c.config.WatchAs("device_state.t_btm", "client_state.t_btm")
	c.config.WatchAs("device_state.power", "client_state.power")
	c.config.WatchAs("device_state.power_diff", "client_state.power_diff")
	c.config.WatchAs("device_state.collect", "client_state.collect")
	c.config.WatchAs("device_state.press", "client_state.press")
	c.config.WatchAs("session_state.error", "client_state.error")
	c.config.WatchAs("session_state.device_id", "client_state.device_id")

	return nil
}

func (c *calc) HandleMessage(_ gen.PID, changeReport any) error {
	return c.config.HandleReport(changeReport)
}
