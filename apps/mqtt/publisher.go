package mqtt

import (
	"samovar/common/models"
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"time"
)

type entry = models.MqttMessage
type queue = map[string]entry

type publisher struct {
	act.Actor
	queue queue
}

func newPublisher() gen.ProcessBehavior {
	return &publisher{}
}

func (p *publisher) Init(_ ...any) error {
	inter.RegisterActor(p, "[mqtt.publisher]")
	p.queue = make(queue)
	return p.schedulePublication()
}

func (p *publisher) HandleMessage(_ gen.PID, msg any) error {
	if entry, ok := msg.(entry); ok {
		p.queue[entry.Topic] = entry
		return nil
	} else {
		return p.publishAll()
	}
}

func (p *publisher) publishAll() error {
	for key, entry := range p.queue {
		if _, err := inter.Call(p, entry, "publication", "mqtt_client"); err == nil {
			delete(p.queue, key)
		} else {
			p.Log().Error("mqtt.publisher.publishAll.error: %s", err)
		}
	}
	return p.schedulePublication()
}

func (p *publisher) schedulePublication() error {
	_, err := p.SendAfter(p.PID(), "publishAll", time.Second)
	return err
}
