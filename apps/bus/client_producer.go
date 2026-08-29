package bus

import (
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
)

type clientProducer struct{ stage.Producer }

func newClientProducer() gen.ProcessBehavior {
	return &clientProducer{}
}

func (p *clientProducer) Init(_ ...any) error {
	p.InitProducer("(bus.client.producer)")

	p.AddReportRoute(
		"client_state",
		"client_state_changed",
	)
	p.AddReportRoute(
		"client_desired_state",
		"client_desired_state_changed",
	)

	return nil
}

func (p *clientProducer) HandleMessage(_ gen.PID, msg any) error {
	return p.HandleChangeReports(msg)
}
