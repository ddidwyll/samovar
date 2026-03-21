package bus

import (
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
	case MqttNewMessage:
		c.Log().Info("bus.deviceConsumer reveive MqttNewMessage: %s", m.Topic)
		return nil
	default:
		err := fmt.Sprintf("bus.mqttProducer receive unexpected event: %#v", event)
		return errors.New(err)
	}
}
