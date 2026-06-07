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
	c.config = clc.NewConfig(c, "{device.calc}")

	c.config.WatchAs("device_raw_state.term_d", "device_state.t_top")
	c.config.WatchAs("device_raw_state.term_c", "device_state.t_mid")
	c.config.WatchAs("device_raw_state.term_k", "device_state.t_btm")
	c.config.WatchAs("device_raw_state.press_a", "device_state.press")
	c.config.WatchAs("device_raw_state.power_m", "device_state.power")

	c.config.Watch(
		calcPowerDiff,
		"device_raw_state.power_m",
		"device_raw_state.power",
	)

	c.config.Watch(
		calcCollect,
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
