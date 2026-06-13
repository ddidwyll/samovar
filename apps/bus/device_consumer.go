package bus

import (
	"samovar/common/models"
	"samovar/lib/change"
	"samovar/lib/inter"
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
)

type deviceConsumer struct{ stage.Consumer }

func newDeviceConsumer() gen.ProcessBehavior { return &deviceConsumer{} }

func (dc *deviceConsumer) Init(_ ...any) error {
	dc.InitConsumer("([bus.device.consumer])")
	dc.SubscribeToEvent("mqtt_new_message")

	dc.AddReportRoute(
		"device_raw_state_changed",
		"device_producer",
		"device_calc",
	)
	dc.AddReportRoute(
		"device_state_changed",
		"device_producer",
		"device_change_log",
	)

	return nil
}

func (dc *deviceConsumer) HandleEvent(event gen.MessageEvent) error {
	if m, ok := event.Message.(models.MqttMessage); ok {
		inter.Trigger(dc, event.Event.Name, "mqtt_producer")
		request := change.NewRequest(m.Topic, m.Text, "mqtt_message", m.Timestamp)
		dc.Log().Debug("bus.deviceConsumer new change.Request: %+v", request)
		return inter.Send(dc, request, "request", "device_raw_state")
	} else {
		return dc.HandleChangeReports(event)
	}
}
