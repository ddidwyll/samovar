package bus

import (
	"samovar/common/models"
	"samovar/lib/change"
	"samovar/lib/stage"

	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type deviceConsumer struct{ stage.Consumer }

func newDeviceConsumer() gen.ProcessBehavior { return &deviceConsumer{} }

func (dc *deviceConsumer) Init(_ ...any) error {
	dc.Log().Debug("bus.deviceConsumer started (%s)", dc.Name())
	return dc.LinkEvents("mqtt_new_message", "device_raw_state_changed")
}

func (dc *deviceConsumer) HandleEvent(event gen.MessageEvent) error {
	switch m := event.Message.(type) {
	case models.MqttMessage:
		request := change.NewRequest(m.Topic, m.Text, "mqtt_message", m.Timestamp)
		dc.Log().Debug("bus.deviceConsumer new change.Request: %+v", request)
		return dc.Send("device_raw_state", request)
	case change.Report:
		dc.Log().Debug("bus.deviceConsumer receive change.Report: %+v", m)
		switch m.LastFrom() {
		case "device_state":
			return dc.Send("device_change_log", m)
		case "device_raw_state":
			return dc.Send("device_state", m.NewRequest())
		default:
			return nil
		}
	default:
		err := fmt.Sprintf("bus.deviceConsumer receive unexpected event: %#v", event)
		return errors.New(err)
	}
}
