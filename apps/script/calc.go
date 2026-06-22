package script

import (
	clc "samovar/lib/calc"

	"ergo.services/ergo/gen"
)

type calc struct{ clc.CalcActor }

var scripts = scriptMap{
	"idle":             performIdle,
	"collect_body_emu": performCollectBodyEmu,
}

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{script.calc}")

	c.Watch(selectScript, "client_desired_state.script")

	c.Watch(
		performScript,
		"script_state.current_script",
	)

	return nil
}

func selectScript(args args, apply applyFn) {
	script := args["client_desired_state.script"].String()
	if _, has := scripts[script]; !has {
		script = "idle"
	}
	apply("script_state.current_script", script)
}

func performScript(args args, apply applyFn) {
	currentScript := args["script_state.current_script"].String()
	for scriptName, scriptFn := range scripts {
		if scriptName == currentScript {
			apply = func(key string, value any) {
				apply("script_state."+key, value)
			}
			scriptFn(args, apply)
		}
	}
}
