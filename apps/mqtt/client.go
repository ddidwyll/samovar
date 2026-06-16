package mqtt

import (
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	natiu "github.com/soypat/natiu-mqtt"

	"errors"
)

type client struct{ act.Actor }

func newClient() gen.ProcessBehavior { return &client{} }

func (s *client) Init(args ...any) error {
	inter.RegisterActor(s, "[mqtt.client]")

	cfg, ok := args[0].(*config)

	if !ok {
		return errors.New("invalid mqtt client config")
	}

	listener := newListener(cfg)

	if _, err := s.SpawnMeta(listener, gen.MetaOptions{}); err != nil {
		return err
	}

	s.Log().Debug("mqtt.client started (%s)", s.Name())

	return nil
}

func (s *client) HandleMessage(_ gen.PID, msg any) error {
	s.Log().Debug("mqtt.client receive message: %#v", msg)
	return inter.Send(s, msg, "message", "mqtt_producer")
}
