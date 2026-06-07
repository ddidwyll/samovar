package bus

import (
	"samovar/lib/change"
	"samovar/lib/inter"
	"samovar/lib/stage"

	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type clientProducer struct{ stage.Producer }

func newClientProducer() gen.ProcessBehavior { return &clientProducer{} }

func (p *clientProducer) Init(_ ...any) error {
	inter.RegisterActor(p, "(bus.client.producer)")
	p.Log().Debug("bus.clientProducer started (%s)", p.Name())

	return p.RegisterEvents("client_state_changed")
}

func (p *clientProducer) HandleMessage(_ gen.PID, msg any) error {
	if report, ok := msg.(change.Report); ok {
		switch report.LastFrom() {
		case "client_state":
			return p.FireEvent("client_state_changed", report)
		default:
			err := fmt.Sprintf("bus.deviceProducer unexpected change.Report: %s", report.LastFrom())
			return errors.New(err)
		}
	} else {
		err := fmt.Sprintf("bus.clientProducer unexpected message: %#v", msg)
		return errors.New(err)
	}
}
