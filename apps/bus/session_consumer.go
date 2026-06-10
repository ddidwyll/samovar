package bus

import (
	"samovar/lib/change"
	"samovar/lib/inter"
	"samovar/lib/stage"

	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type sessionConsumer struct{ stage.Consumer }

func newSessionConsumer() gen.ProcessBehavior {
	return &sessionConsumer{}
}

func (sc *sessionConsumer) Init(_ ...any) error {
	inter.RegisterActor(sc, "([bus.session.consumer])")
	sc.Log().Debug("bus.sessionConsumer started (%s)", sc.Name())

	return sc.LinkEvents(
		"device_state_changed",
		"session_state_changed",
	)
}

func (sc *sessionConsumer) HandleEvent(event gen.MessageEvent) error {
	if report, ok := event.Message.(change.Report); ok {
		sc.Log().Debug("bus.sessionConsumer <- change.Report[%s]", report.LastFrom())

		switch report.LastFrom() {
		case "device_state":
			inter.Trigger(sc, event.Event.Name, "device_producer")
			return inter.Send(sc, report, "report", "session_calc")
		case "session_state":
			inter.Trigger(sc, event.Event.Name, "session_producer")
			return inter.Send(sc, report, "report", "session_calc")
		default:
			return nil
		}
	} else {
		err := fmt.Sprintf("bus.sessionConsumer receive unexpected event: %#v", event)
		return errors.New(err)
	}
}
