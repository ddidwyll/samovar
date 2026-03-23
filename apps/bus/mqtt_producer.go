package bus

import (
	"samovar/common/models"
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
	"errors"
	"fmt"
)

type mqttProducer struct{ stage.Producer }

func newMqttProducer() gen.ProcessBehavior { return &mqttProducer{} }

func (p *mqttProducer) Init(_ ...any) error {
	p.Log().Debug("bus.mqttProducer started (%s)", p.Name())
	return p.RegisterEvents("mqtt_new_message")
}

func (p *mqttProducer) HandleMessage(_ gen.PID, msg any) error {
	switch message := msg.(type) {
	case models.MqttMessage:
		p.Log().Debug("bus.mqttProducer receive MqttMessage: %v", message)
		return p.FireEvent("mqtt_new_message", message)
	default:
		err := fmt.Sprintf("bus.mqttProducer receive unexpected message: %#v", msg)
		return errors.New(err)
	}
}
