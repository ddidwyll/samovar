package bus

import (
	"samovar/common/models"
	"samovar/lib/change"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"errors"
	"fmt"
)

var deviceEvents = []gen.Atom{
	"mqtt_new_message",
}

type deviceConsumer struct {
	act.Actor
}

func newDeviceConsumer() gen.ProcessBehavior {
	return &deviceConsumer{}
}

func (c *deviceConsumer) Init(_ ...any) error {
	nodeName := c.Node().Name()

	for _, e := range deviceEvents {
		if _, err := c.LinkEvent(gen.Event{e, nodeName}); err != nil {
			return err
		}
	}

	c.Log().Debug("bus.deviceConsumer started (%s)", c.Name())
	return nil
}

func (c *deviceConsumer) HandleEvent(event gen.MessageEvent) error {
	switch m := event.Message.(type) {
	case models.MqttMessage:
		ch := change.Request{m.Topic, m.Text, m.Timestamp}
		c.Log().Debug("bus.deviceConsumer new ChangeRequest request: %+v", ch)
		return c.Send("device_raw_state", ch)
	default:
		err := fmt.Sprintf("bus.mqttProducer receive unexpected event: %#v", event)
		return errors.New(err)
	}
}
