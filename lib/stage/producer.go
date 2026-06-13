package stage

import (
	"samovar/lib/change"
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
)

const (
	notify = false
	buffer = 100
)

type eventRefs map[gen.Atom]gen.Ref
type producerRoutes map[string]gen.Atom

type Producer struct {
	act.Actor
	refs   eventRefs
	routes producerRoutes
}

func (p *Producer) InitProducer(name string) {
	inter.RegisterActor(p, name)
	p.refs = make(eventRefs)
	p.routes = make(producerRoutes)
}

func (p *Producer) AddReportRoute(from string, event gen.Atom) {
	if err := p.RegisterEvents(event); err == nil {
		p.routes[from] = event
	} else {
		panic(err)
	}
}

func (p *Producer) RegisterEvents(events ...gen.Atom) error {
	opts := gen.EventOptions{notify, buffer}

	for _, event := range events {
		if _, exists := p.refs[event]; exists {
			continue
		}
		if ref, err := p.RegisterEvent(event, opts); err == nil {
			p.refs[event] = ref
		} else {
			return err
		}
	}

	return nil
}

func (p *Producer) FireEvent(event gen.Atom, msg any) error {
	ref := p.refs[event]
	return p.SendEvent(event, ref, msg)
}

func (p *Producer) HandleChangeReports(msg any) error {
	switch r := msg.(type) {
	case []change.Report:
		return p.routeReports(r...)
	case change.Report:
		return p.routeReports(r)
	default:
		return errors.New("Unexpected producer message")
	}
}

func (p *Producer) routeReports(reports ...change.Report) error {
	for _, report := range reports {
		if event, exists := p.routes[report.LastFrom()]; !exists {
			return errors.New("Unexpected producer message")
		} else {
			if err := p.FireEvent(event, report); err != nil {
				return err
			}
		}
	}
	return nil
}
