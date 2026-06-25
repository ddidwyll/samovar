package client

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/gen"
)

type calc struct{ clc.CalcActor }

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{client.calc}")

	// DEVICE_STATE
	c.WatchAs("device_state.t_top", "client_state.t_top")
	c.WatchAs("device_state.t_mid", "client_state.t_mid")
	c.WatchAs("device_state.t_btm", "client_state.t_btm")
	c.WatchAs("device_state.power", "client_state.power")
	c.WatchAs("device_state.power_diff", "client_state.power_diff")
	c.WatchAs("device_state.collect", "client_state.collect")
	c.WatchAs("device_state.press", "client_state.press")
	// DEVICE_RAW_STATE
	c.WatchAs("device_raw_state.last_rx", "client_state.last_mqtt_rx")
	c.WatchAs("device_raw_state.last_tx", "client_state.last_mqtt_tx")
	c.WatchAs("device_raw_state.last_ping", "client_state.last_mqtt_ping")
	// SESSION_STATE
	c.WatchAs("session_state.error", "client_state.error")
	c.WatchAs("session_state.started", "client_state.started")
	c.WatchAs("session_state.device_id", "client_state.device_id")
	c.WatchAs("session_state.devices", "client_state.devices")
	// SCRIPT_STATE
	c.WatchAs("script_state.scripts", "client_state.scripts")
	c.WatchAs("script_state.script_mode", "client_state.script_mode")

	return nil
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	return c.HandleChangeReports(msg)
}
