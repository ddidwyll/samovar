package mqtt

import (
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
)

type client struct {
	act.Actor
	config   *config
	listener gen.Alias
}

func newClient() gen.ProcessBehavior { return &client{} }

func (c *client) Init(args ...any) error {
	inter.RegisterActor(c, "[mqtt.client]")

	if err := c.initConfig(args...); err != nil {
		return err
	}

	return c.createListener()
}

func (c *client) initConfig(args ...any) error {
	if cfg, ok := args[0].(*config); !ok {
		return errors.New("invalid mqtt client config")
	} else {
		c.config = cfg
		return nil
	}
}

func (c *client) createListener() error {
	listener := newListener(c.config)
	opts := gen.MetaOptions{}
	if alias, err := c.SpawnMeta(listener, opts); err != nil {
		return err
	} else {
		c.listener = alias
		return nil
	}
}

func (c *client) HandleMessage(_ gen.PID, msg any) error {
	c.Log().Debug("mqtt.client receive message: %#v", msg)
	return inter.Send(c, msg, "message", "mqtt_producer")
}

func (c *client) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	c.Log().Info("mqtt.client receive request: %#v", req)
	return "", nil
}
