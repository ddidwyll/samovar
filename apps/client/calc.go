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

	c.config.Watch(
		clc.TryAsIs,
		"device_state.t_top",
		"client_state.t_top",
	)

	return nil
}

func (c *calc) HandleMessage(_ gen.PID, changeReport any) error {
	return c.config.HandleReport(changeReport)
}
