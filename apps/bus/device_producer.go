package bus

import (
	"samovar/lib/change"
	"samovar/lib/stage"

	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type deviceProducer struct{ stage.Producer }

func newDeviceProducer() gen.ProcessBehavior { return &deviceProducer{} }

func (p *deviceProducer) Init(_ ...any) error {
	p.Log().Debug("bus.deviceProducer started (%s)", p.Name())

	return p.RegisterEvents(
		"device_raw_state_changed",
		"device_state_changed",
	)
}

func (p *deviceProducer) HandleMessage(_ gen.PID, msg any) error {
	if report, ok := msg.(change.Report); ok {
		p.Log().Debug("bus.deviceProducer received change.Report: %v", report)

		switch report.LastFrom() {
		case "device_raw_state":
			return p.FireEvent("device_raw_state_changed", report)
		case "device_state":
			return p.FireEvent("device_state_changed", report)
		default:
			err := fmt.Sprintf("bus.deviceProducer unexpected change.Report: %s", report.LastFrom())
			return errors.New(err)
		}
	} else {
		err := fmt.Sprintf("bus.deviceProducer unexpected message: %#v", msg)
		return errors.New(err)
	}
}
