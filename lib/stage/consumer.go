package stage

import (
	"samovar/lib/change"
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
)

type route struct {
	recipient string
	producer  gen.Atom
}

type routes map[gen.Atom]route

type Consumer struct {
	act.Actor
	routes routes
}

func (c *Consumer) InitConsumer(name string) {
	inter.RegisterActor(c, name)
	c.routes = make(routes)
}

func (c *Consumer) AddRoute(event, producer gen.Atom, recipient string) {
	if err := c.Subscribe(event); err == nil {
		c.routes[event] = route{recipient, producer}
	} else {
		panic(err)
	}
}

func (c *Consumer) Subscribe(e gen.Atom) error {
	nodeName := c.Node().Name()
	event := gen.Event{e, nodeName}
	_, err := c.LinkEvent(event)
	return err
}

func (c *Consumer) HandleReports(e gen.MessageEvent) error {
	eventName := e.Event.Name
	switch r := e.Message.(type) {
	case []change.Report:
		return c.routeReports(eventName, r)
	case change.Report:
		return c.routeReports(eventName, []change.Report{r})
	default:
		return errors.New("Unexpected consumer event")
	}
}

func (c *Consumer) routeReports(event gen.Atom, reports []change.Report) error {
	for _, report := range reports {
		if route, exists := c.routes[event]; !exists {
			return errors.New("Unexpected consumer event")
		} else {
			inter.Trigger(c, event, route.producer)
			if err := inter.Send(c, report, "report", route.recipient); err != nil {
				return err
			}
		}
	}
	return nil
}
