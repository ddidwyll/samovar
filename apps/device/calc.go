package device

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
	c.config = clc.NewConfig(c)

	c.config.Watch(
		clc.TryAsIs,
		"device_state.t_top",
		"device_raw_state.term_d",
	)

	c.config.Watch(
		clc.TryAsIs,
		"device_state.t_mid",
		"device_raw_state.term_c",
	)

	c.config.Watch(
		clc.TryAsIs,
		"device_state.t_btm",
		"device_raw_state.term_k",
	)

	c.config.Watch(
		clc.TryAsIs,
		"device_state.press",
		"device_raw_state.press_a",
	)

	c.config.Watch(
		clc.TryAsIs,
		"device_state.power",
		"device_raw_state.power_m",
	)

	c.config.Watch(
		calcPowerDiff,
		"device_state.power_diff",
		"device_raw_state.power_m",
		"device_raw_state.power",
	)

	c.config.Watch(
		calcCollect,
		"device_state.collect",
		"device_raw_state.otbor",
		"device_raw_state.flag_otb",
		"device_raw_state.sek_otb",
		"device_raw_state.min_otb",
	)

	return nil
}

func (c *calc) HandleMessage(_ gen.PID, changeReport any) error {
	return c.config.HandleReport(changeReport)
}
