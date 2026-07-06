package device

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/gen"

	"time"
)

type calc struct{ clc.CalcActor }

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{device.calc}")

	// DEVICE_RAW_STATE
	c.WatchReport(tick, "device_raw_state.last_ping")

	c.WatchFields(calcDecigrad, "device_raw_state.term_d")
	c.WatchFields(calcDecigrad, "device_raw_state.term_c")
	c.WatchFields(calcDecigrad, "device_raw_state.term_k")
	c.WatchFieldAs("device_raw_state.press_a", "device_state.press")
	c.WatchFieldAs("device_raw_state.power_m", "device_state.power")

	c.WatchFields(
		calcPowerDiff,
		"device_raw_state.power_m",
		"device_raw_state.power",
	)
	c.WatchFields(
		calcCollect,
		"device_raw_state.otbor",
		"device_raw_state.flag_otb",
		"device_raw_state.sek_otb",
		"device_raw_state.min_otb",
		"device_raw_state.tick",
	)

	// SESSION_DESIRED_STATE
	c.WatchFieldAs(
		"session_desired_state.power",
		"device_desired_state.power",
	)

	// SCRIPT_STATE
	c.WatchFieldAs(
		"script_state.collect_type",
		"device_desired_state.collect_type",
	)
	c.WatchFieldAs(
		"script_state.collect_value",
		"device_desired_state.collect_value",
	)

	// DEVICE_STATE
	//

	return nil
}

func tick(r clc.Report, _ clc.FetchFn, apply clc.ApplyFn) {
	secA := time.UnixMilli(r.OldTimestamp).Second()
	secB := time.UnixMilli(r.NewTimestamp).Second()
	if secA != secB {
		apply("device_raw_state.tick", r.NewTimestamp)
	}
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	return c.HandleChangeReports(msg)
}
