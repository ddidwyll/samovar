package bus

import (
	"samovar/lib/change"
	"samovar/lib/inter"
	"samovar/lib/stage"

	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type clientConsumer struct{ stage.Consumer }

func newClientConsumer() gen.ProcessBehavior {
	return &clientConsumer{}
}

func (cc *clientConsumer) Init(_ ...any) error {
	inter.RegisterActor(cc, "([bus.client.consumer])")
	cc.Log().Debug("bus.clientConsumer started (%s)", cc.Name())

	return cc.LinkEvents(
		"device_state_changed",
		"client_state_changed",
		"session_state_changed",
	)
}

func (cc *clientConsumer) HandleEvent(event gen.MessageEvent) error {
	if report, ok := event.Message.(change.Report); ok {
		cc.Log().Debug("bus.clientConsumer <- change.Report[%s]", report.LastFrom())

		switch report.LastFrom() {
		case "device_state":
			inter.Trigger(cc, event.Event.Name, "device_producer")
			return inter.Send(cc, report, "report", "client_calc")
		case "session_state":
			inter.Trigger(cc, event.Event.Name, "session_producer")
			return inter.Send(cc, report, "report", "client_calc")
		case "client_state":
			inter.Trigger(cc, event.Event.Name, "client_producer")
			return inter.Send(cc, report, "report", "client_feed")
		default:
			return nil
		}
	} else {
		err := fmt.Sprintf("bus.clientConsumer receive unexpected event: %#v", event)
		return errors.New(err)
	}
}
