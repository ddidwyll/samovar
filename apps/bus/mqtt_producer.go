package bus

import (
	"samovar/apps/mqtt"

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
	opts := gen.EventOptions{true, 100}

	for _, e := range MqttEvents {
		if ref, err := p.RegisterEvent(e, opts); err == nil {
			p.eventRefs[e] = ref
		} else {
			return err
		}
	}

	return nil
}

func (p *mqttProducer) HandleMessage(from gen.PID, msg any) error {
	switch m := msg.(type) {
	case mqtt.Message:
		event := gen.Atom("mqtt_new_message")
		ref := p.eventRefs[event]
		p.Log().Debug("bus.mqttProducer receive MqttNewMessage: %v", m)
		return p.SendEvent(event, ref, MqttNewMessage{m})
	default:
		err := fmt.Sprintf("bus.mqttProducer receive unexpected message: %#v", msg)
		return errors.New(err)
	}
}
