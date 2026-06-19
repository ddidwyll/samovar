package bus

import (
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
)

type deviceProducer struct{ stage.Producer }

func newDeviceProducer() gen.ProcessBehavior {
	return &deviceProducer{}
}

func (p *deviceProducer) Init(_ ...any) error {
	p.InitProducer("(bus.device.producer)")

	p.AddReportRoute(
		"device_raw_state",
		"device_raw_state_changed",
	)
	p.AddReportRoute(
		"device_desired_state",
		"device_desired_state_changed",
	)
	p.AddReportRoute(
		"device_state",
		"device_state_changed",
	)

	return nil
}

func (p *deviceProducer) HandleMessage(_ gen.PID, msg any) error {
	return p.HandleChangeReports(msg)
}
