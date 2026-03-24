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

func (c *deviceConsumer) Init(_ ...any) error {
	c.Log().Debug("bus.deviceConsumer started (%s)", c.Name())
	return c.LinkEvents("mqtt_new_message", "device_raw_state_changed")
}

func (c *deviceConsumer) HandleEvent(event gen.MessageEvent) error {
	switch m := event.Message.(type) {
	case models.MqttMessage:
		request := change.NewRequest(m.Topic, m.Text, "mqtt_message", m.Timestamp)
		c.Log().Debug("bus.deviceConsumer new change.Request: %+v", request)
		return c.Send("device_raw_state", request)
  case change.Report:
		c.Log().Debug("bus.deviceConsumer receive change.Report: %+v", m)
		return c.Send("device_change_log", m)
	default:
		err := fmt.Sprintf("bus.mqttProducer receive unexpected event: %#v", event)
		return errors.New(err)
	}
}
