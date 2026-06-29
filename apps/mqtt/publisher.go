package mqtt

import (
	"samovar/common/models"
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"fmt"
	"time"
)

type queueEntry = models.MqttMessage
type queue = map[string]queueEntry

type debounceEntry struct {
	count int64
	value string
	ts    int64
}

type debounce = map[string]debounceEntry

type publisher struct {
	act.Actor
	queue    queue
	debounce debounce
}

var publishOrder = []string{
  "otbor_t_new",
  "otbor_g_1_new",
  "otbor_g_2_new",
  "power_m_new",
  "work",
  "otbor_new",
}

func newPublisher() gen.ProcessBehavior {
	return &publisher{}
}

func (p *publisher) Init(_ ...any) error {
	inter.RegisterActor(p, "[mqtt.publisher]")
	p.queue = make(queue)
	p.debounce = make(debounce)
	return p.schedulePublication()
}

func (p *publisher) HandleMessage(_ gen.PID, msg any) error {
	if entry, ok := msg.(queueEntry); ok {
		p.queue[entry.Topic] = entry
		return nil
	} else {
		return p.publishAll()
	}
}

func (p *publisher) publishAll() error {
	for _, key := range publishOrder {
		entry, exists := p.queue[key]
		if !exists || p.skipDebounce(entry) {
			continue
		}
		if _, err := inter.Call(p, entry, "publication", "mqtt_client"); err == nil {
			delete(p.queue, key)
			p.checkDebounce(entry)
			time.Sleep(time.Second)
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

func (p *publisher) skipDebounce(qe queueEntry) bool {
	now := time.Now().Unix()
	de, exists := p.debounce[qe.Topic]
	if !exists || de.value != qe.Text {
		return false
	}
	diff := now - de.ts
	return diff < de.count*5
}

func (p *publisher) checkDebounce(qe queueEntry) {
	now := time.Now().Unix()
	de, exists := p.debounce[qe.Topic]
	if !exists || de.value != qe.Text {
		fmt.Printf("##### mqtt.publisher.debounce.check: %s, 0, %s\n", qe.Topic, qe.Text)
		p.debounce[qe.Topic] = debounceEntry{1, qe.Text, now}
	} else {
		p.debounce[qe.Topic] = debounceEntry{de.count + 1, de.value, now}
		fmt.Printf("##### mqtt.publisher.debounce.check: %s, %d, %s\n", qe.Topic, de.count+1, qe.Text)
	}
}
