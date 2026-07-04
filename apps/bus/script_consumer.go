package bus

import (
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
)

type scriptConsumer struct{ stage.Consumer }

func newScriptConsumer() gen.ProcessBehavior {
	return &scriptConsumer{}
}

func (sc *scriptConsumer) Init(_ ...any) error {
	sc.InitConsumer("([bus.script.consumer])")

	sc.AddReportRoute(
		"client_desired_state_changed",
		"script_calc",
	)

	sc.AddReportRoute(
		"script_state_changed",
		"script_calc",
	)

	sc.AddReportRoute(
		"session_state_changed",
		"script_calc",
	)

	sc.Log().Debug("bus.scriptConsumer.Init.name: %s", sc.Name())
	return nil
}

func (sc *scriptConsumer) HandleEvent(event gen.MessageEvent) error {
	return sc.HandleChangeReports(event)
}
