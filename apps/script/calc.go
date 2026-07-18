package script

import (
	clc "samovar/lib/calc"
	ss "samovar/lib/scripts"

	"ergo.services/ergo/gen"

	"fmt"
	"maps"
	"slices"
	"strings"
)

type calc struct{ clc.CalcActor }

var scripts = scriptMap{
	"idle":             performIdle,
	"collect_head_min": ss.PerformCollectHeadMin,
	"collect_waste":    ss.PerformCollectWaste,
	"collect_recyc_slow":    ss.PerformCollectRecycSlow,
	"collect_recyc_fast":    ss.PerformCollectRecycFast,
	"collect_body":     ss.PerformCollectBody,
}

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{script.calc}")

	c.WatchFields(
		selectScript,
		"client_desired_state.script_mode",
	)

	c.WatchFields(
		performScript,
		"script_state.script_mode",
		"device_state.t_btm",
		"session_state.min_reflux_ratio",
		"session_state.net_power",
		"session_state.max_collect",
		"session_state.max_stable_temp_diff",
		"session_state.mid_stable_diff",
		"session_state.mid_stable_temp",
		"session_state.is_mid_stable",
	)

	return c.calcScripts()
}

func (c *calc) calcScripts() error {
	modes := slices.Collect(maps.Keys(scripts))
	c.Log().Debug("script.calc.calcScripts.modes: %v", modes)
	return c.SendRequest("script_state.scripts", modes)
}

func selectScript(args args, apply applyFn) {
	script := args.MustGet("client_desired_state.script_mode").String()
	if _, has := scripts[script]; !has {
		script = "idle"
	}
	apply("script_state.script_mode", script)
}

func performScript(args args, apply applyFn) {
	scriptMode := args.MustGet("script_state.script_mode").String()
	for scriptName, scriptFn := range scripts {
		if scriptName == scriptMode {
			patchedApply := func(key string, value any) {
				if !strings.Contains(key, ".") {
					key = "script_state." + key
				}
				fmt.Printf(">>> script.calc.performScript.patchedApply: %v, %v\n", key, value)
				apply(key, value)
			}
			scriptFn(args, patchedApply)
		}
	}
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	return c.HandleChangeReports(msg)
}
