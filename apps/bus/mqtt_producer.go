package bus

import (
	"samovar/common/models"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"errors"
	"fmt"
)

var MqttEvents = []gen.Atom{
	"mqtt_new_message",
}

type mqttProducer struct {
	act.Actor
	eventRefs map[gen.Atom]gen.Ref
}

func newMqttProducer() gen.ProcessBehavior {
	return &mqttProducer{eventRefs: make(map[gen.Atom]gen.Ref)}
}

func (p *mqttProducer) Init(_ ...any) error {
	opts := gen.EventOptions{false, 0}

	for _, e := range MqttEvents {
		if ref, err := p.RegisterEvent(e, opts); err == nil {
			p.eventRefs[e] = ref
		} else {
			return err
		}
	}

	p.Log().Debug("bus.mqttProducer started (%s)", p.Name())
	return nil
}

func (p *mqttProducer) HandleMessage(_ gen.PID, msg any) error {
	switch m := msg.(type) {
	case models.MqttMessage:
		event := gen.Atom("mqtt_new_message")
		ref := p.eventRefs[event]
		p.Log().Debug("bus.mqttProducer receive MqttMessage: %v", m)
		return p.SendEvent(event, ref, m)
	default:
		err := fmt.Sprintf("bus.mqttProducer receive unexpected message: %#v", msg)
		return errors.New(err)
	}
}
