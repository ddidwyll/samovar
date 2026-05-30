package bus

import (
	"samovar/lib/change"
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
	cc.Log().Debug("bus.clientConsumer started (%s)", cc.Name())

	return cc.LinkEvents(
		"device_state_changed",
		"client_state_changed",
	)
}

func (cc *clientConsumer) HandleEvent(event gen.MessageEvent) error {
	if report, ok := event.Message.(change.Report); ok {
		cc.Log().Debug("bus.clientConsumer <- change.Report[%s]", report.LastFrom())

		switch report.LastFrom() {
		case "device_state":
			return cc.Send("client_state", report.ReRequest())
		case "client_state":
			return cc.Send("client_feed", report)
		default:
			return nil
		}
	} else {
		err := fmt.Sprintf("bus.clientConsumer receive unexpected event: %#v", event)
		return errors.New(err)
	}
}
