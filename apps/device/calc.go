package device

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/gen"
)

type calc struct{ clc.CalcActor }

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{device.calc}")

	// DEVICE_RAW_STATE
	c.WatchAs("device_raw_state.term_d", "device_state.t_top")
	c.WatchAs("device_raw_state.term_c", "device_state.t_mid")
	c.WatchAs("device_raw_state.term_k", "device_state.t_btm")
	c.WatchAs("device_raw_state.press_a", "device_state.press")
	c.WatchAs("device_raw_state.power_m", "device_state.power")

	c.Watch(
		calcPowerDiff,
		"device_raw_state.power_m",
		"device_raw_state.power",
	)

	c.Watch(
		calcCollect,
		"device_raw_state.otbor",
		"device_raw_state.flag_otb",
		"device_raw_state.sek_otb",
		"device_raw_state.min_otb",
	)

	// SESSION_DESIRED_STATE
	c.WatchAs(
		"session_desired_state.power",
		"device_desired_state.power",
	)

	// SCRIPT_STATE
	c.WatchAs(
		"script_state.collect_type",
		"device_desired_state.collect_type",
	)
	c.WatchAs(
		"script_state.collect_value",
		"device_desired_state.collect_value",
	)

	return nil
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	return c.HandleChangeReports(msg)
}
