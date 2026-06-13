package bus

import (
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
)

type sessionConsumer struct{ stage.Consumer }

func newSessionConsumer() gen.ProcessBehavior {
	return &sessionConsumer{}
}

func (sc *sessionConsumer) Init(_ ...any) error {
	sc.InitConsumer("([bus.session.consumer])")

	sc.AddReportRoute(
		"device_state_changed",
		"device_producer",
		"session_calc",
	)
	sc.AddReportRoute(
		"session_state_changed",
		"session_producer",
		"session_calc",
	)

	sc.Log().Debug("bus.sessionConsumer started (%s)", sc.Name())
	return nil
}

func (sc *sessionConsumer) HandleEvent(event gen.MessageEvent) error {
	return sc.HandleChangeReports(event)
}
