package client

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/gen"

	"fmt"
	// "math"
)

type calc struct{ clc.CalcActor }

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{client.calc}")

	// DEVICE_STATE
	c.WatchFields(calcDegree, "device_state.t_top")
	c.WatchFields(calcDegree, "device_state.t_mid")
	c.WatchFields(calcDegree, "device_state.t_btm")
	c.WatchFieldAs("device_state.power", "client_state.power")
	c.WatchFieldAs("device_state.power_diff", "client_state.power_diff")
	c.WatchFieldAs("device_state.collect", "client_state.collect")
	c.WatchFieldAs("device_state.press", "client_state.press")
	// DEVICE_RAW_STATE
	// c.WatchFieldAs("device_raw_state.last_rx", "client_state.last_mqtt_rx")
	// c.WatchFieldAs("device_raw_state.last_tx", "client_state.last_mqtt_tx")
	c.WatchFieldAs("device_raw_state.last_ping", "client_state.last_mqtt_ping")
	c.WatchFieldAs("device_raw_state.collect_synced", "client_state.collect_synced")
	// SESSION_STATE
	c.WatchFieldAs("session_state.error", "client_state.error")
	c.WatchFieldAs("session_state.ready", "client_state.ready")
	c.WatchFieldAs("session_state.device_id", "client_state.device_id")
	c.WatchFieldAs("session_state.devices", "client_state.devices")
	c.WatchFieldAs("session_state.collection_speed", "client_state.collection_speed")
	c.WatchFieldAs("session_state.is_mid_stable", "client_state.is_mid_stable")
	c.WatchFields(calcDegree, "session_state.mid_stable_temp")
	c.WatchFields(calcDegree, "session_state.mid_stable_diff")
	c.WatchFieldAs("session_state.net_power", "client_state.net_power")
	c.WatchFieldAs("session_state.reflux_ratio", "client_state.reflux_ratio")
	c.WatchFields(
		calcCollected,
		"session_state.body_collected_value",
		"session_state.head_collected_value",
		"session_state.recyc_collected_value",
		"session_state.waste_collected_value",
	)
	c.WatchFields(
		calcAverageCollectSpeed,
		"session_state.body_collect_duration",
		"session_state.head_collect_duration",
		"session_state.recyc_collect_duration",
		"session_state.body_collected_value",
		"session_state.head_collected_value",
		"session_state.recyc_collected_value",
	)
	// SCRIPT_STATE
	c.WatchFieldAs("script_state.scripts", "client_state.scripts")
	c.WatchFieldAs("script_state.script_mode", "client_state.script_mode")

	return nil
}

func calcDegree(args clc.Args, apply clc.ApplyFn) {
	for key, val := range args {
		if !val.IsInt() {
			continue
		}
		target := clc.Fid(key)
		targetKey := "client_state." + target.FieldKey
		apply(targetKey, val.ToFlt()/10.0)
	}
}

func calcAverageCollectSpeed(args clc.Args, apply clc.ApplyFn) {
	accPrefixes := []string{"body", "head", "recyc"}

	for _, prefix := range accPrefixes {
		durationMs := args.MustGet("session_state." + prefix + "_collect_duration")
		collectedMg := args.MustGet("session_state." + prefix + "_collected_value")
		// fmt.Printf(">>> calcAverageCollectSpeed[%s]: dur=%v, col=%v\n", prefix, durationMs, collectedMg)
		acc := "client_state." + prefix + "_average_speed"
		if !durationMs.IsInt() || !collectedMg.IsInt() {
			apply(acc, 0)
			continue
		}
		speedGH := collectedMg.ToInt() * 3600 / durationMs.ToInt()
		// fmt.Printf("calcAverageCollectSpeed.speed[%s]: %dg/h\n", prefix, speedGH)
		apply(acc, speedGH)
	}
}

func calcCollected(args clc.Args, apply clc.ApplyFn) {
	collectAccs := []string{"body_collected", "head_collected", "recyc_collected", "waste_collected"}

	for _, acc := range collectAccs {
		collectedMg := args.MustGet("session_state." + acc + "_value")
		acc := "client_state." + acc
		if !collectedMg.IsInt() {
			apply(acc, 0)
		} else {
			g := collectedMg.ToInt() / 1000
			ml := collectedMg.ToInt() / 789

			apply(acc, fmt.Sprintf("%dg/%dml", g, ml))
		}
	}
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	return c.HandleChangeReports(msg)
}
