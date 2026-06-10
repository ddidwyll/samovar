package session

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
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

	return nil
}

func (c *calc) applyDeviceConfig(args ...any) error {
	cfg, ok := args[0].(*config)
	if !ok {
		return errors.New("session.calc: invalid device config")
	}
}

func (c *calc) HandleMessage(_ gen.PID, changeReport any) error {
	return c.config.HandleReport(changeReport)
}
