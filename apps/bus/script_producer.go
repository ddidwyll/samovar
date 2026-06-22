package bus

import (
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
)

type scriptProducer struct{ stage.Producer }

func newScriptProducer() gen.ProcessBehavior {
	return &scriptProducer{}
}

func (p *scriptProducer) Init(_ ...any) error {
	p.InitProducer("(bus.script.producer)")

	p.AddReportRoute(
		"script_state",
		"script_state_changed",
	)

	return nil
}

func (p *scriptProducer) HandleMessage(_ gen.PID, msg any) error {
	p.Log().Debug("bus.scriptProducer.HandleMessage.msg: %v", msg)
	return p.HandleChangeReports(msg)
}
