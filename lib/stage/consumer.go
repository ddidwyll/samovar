package stage

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Consumer struct {
	act.Actor
}

func (c *Consumer) LinkEvents(events ...gen.Atom) error {
	nodeName := c.Node().Name()

	for _, e := range events {
		event := gen.Event{e, nodeName}
		if _, err := c.LinkEvent(event); err != nil {
			return err
		}
	}

	return nil
}
