package bus

import (
	"samovar/lib/change"
	"samovar/lib/inter"
	"samovar/lib/stage"

	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type sessionProducer struct{ stage.Producer }

func newSessionProducer() gen.ProcessBehavior { return &sessionProducer{} }

func (p *sessionProducer) Init(_ ...any) error {
	inter.RegisterActor(p, "(bus.session.producer)")
	p.Log().Debug("bus.sessionProducer started (%s)", p.Name())

	return p.RegisterEvents("session_state_changed")
}

func (p *sessionProducer) HandleMessage(_ gen.PID, msg any) error {
	if report, ok := msg.(change.Report); ok {
		switch report.LastFrom() {
		case "session_state":
			return p.FireEvent("session_state_changed", report)
		default:
			err := fmt.Sprintf("bus.sessionProducer unexpected change.Report: %s", report.LastFrom())
			return errors.New(err)
		}
	} else {
		err := fmt.Sprintf("bus.sessionProducer unexpected message: %#v", msg)
		return errors.New(err)
	}
}
