package script

import (
	clc "samovar/lib/calc"
	ss "samovar/lib/scripts"

	"ergo.services/ergo/gen"

	"fmt"
	"maps"
	"slices"
)

type calc struct{ clc.CalcActor }

var scripts = scriptMap{
	"idle":               performIdle,
	"collect_head_min":   ss.PerformCollectHeadMin,
	"collect_recyc_fast": ss.PerformCollectRecycFast,
	"collect_recyc_slow": ss.PerformCollectRecycSlow,
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
				fmt.Printf(">>> script.calc.performScript.patchedApply: %v, %v\n", key, value)
				apply("script_state."+key, value)
			}
			scriptFn(args, patchedApply)
		}
	}
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	return c.HandleChangeReports(msg)
}
