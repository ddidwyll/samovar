package bus

import (
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
)

type sessionProducer struct{ stage.Producer }

func newSessionProducer() gen.ProcessBehavior {
	return &sessionProducer{}
}

func (p *sessionProducer) Init(_ ...any) error {
	p.InitProducer("(bus.session.producer)")

	p.AddReportRoute(
		"session_state",
		"session_state_changed",
	)

	return nil
}

func (p *sessionProducer) HandleMessage(_ gen.PID, msg any) error {
	return p.HandleChangeReports(msg)
}
