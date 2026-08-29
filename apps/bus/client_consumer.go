package bus

import (
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
)

type clientConsumer struct{ stage.Consumer }

func newClientConsumer() gen.ProcessBehavior {
	return &clientConsumer{}
}

func (cc *clientConsumer) Init(_ ...any) error {
	cc.InitConsumer("([bus.client.consumer])")

	cc.AddReportRoute(
		"device_raw_state_changed",
		"client_calc",
	)
	cc.AddReportRoute(
		"device_state_changed",
		"client_calc",
	)
	cc.AddReportRoute(
		"session_state_changed",
		"client_calc",
	)
	cc.AddReportRoute(
		"script_state_changed",
		"client_calc",
	)
	cc.AddReportRoute(
		"client_state_changed",
		"client_feed",
	)

	return nil
}

func (cc *clientConsumer) HandleEvent(event gen.MessageEvent) error {
	return cc.HandleChangeReports(event)
}
