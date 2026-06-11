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

	cc.AddRoute(
		"device_state_changed",
		"device_producer",
		"client_calc",
	)

  cc.AddRoute(
		"session_state_changed",
		"session_producer",
		"client_calc",
	)

  cc.AddRoute(
		"client_state_changed",
		"client_producer",
		"client_feed",
	)

	return nil
}

func (cc *clientConsumer) HandleEvent(event gen.MessageEvent) error {
  return cc.HandleReports(event)
}
