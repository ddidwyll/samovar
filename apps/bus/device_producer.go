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
	return p.RegisterEvents("device_raw_state_changed")
}

func (p *deviceProducer) HandleMessage(_ gen.PID, msg any) error {
	report, ok := msg.(change.Report)

	if !ok {
		err := fmt.Sprintf("bus.deviceProducer unexpected message: %#v", msg)
		return errors.New(err)
	} else {
		p.Log().Debug("bus.deviceProducer received change.Report: %v", report)
	}

	switch report.LastFrom() {
	case "raw_state":
		p.FireEvent("device_raw_state_changed", report)
	default:
		p.Log().Error("bus.deviceProducer unexpected change.Report: %s", report.LastFrom())
	}

	return nil
}
