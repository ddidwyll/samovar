package stage

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Producer struct {
	act.Actor
	refs map[gen.Atom]gen.Ref
}

const (
	notify = false
	buffer = 100
)

func (p *Producer) RegisterEvents(events ...gen.Atom) error {
	opts := gen.EventOptions{notify, buffer}
	p.refs = make(map[gen.Atom]gen.Ref)

	for _, event := range events {
		if ref, err := p.RegisterEvent(event, opts); err == nil {
			p.refs[event] = ref
		} else {
			return err
		}
	}

	return nil
}

func (p *Producer) FireEvent(e string, msg any) error {
	event := gen.Atom(e)
	ref := p.refs[event]

	return p.SendEvent(event, ref, msg)
}
