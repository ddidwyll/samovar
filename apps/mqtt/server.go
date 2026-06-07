package mqtt

import (
	"samovar/lib/inter"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
)

type server struct{ act.Actor }

func newServer() gen.ProcessBehavior { return &server{} }

func (s *server) Init(args ...any) error {
	inter.RegisterActor(s, "[mqtt.server]")

	cfg, ok := args[0].(*config)

	if !ok {
		return errors.New("invalid mqtt client config")
	}

	listener := newListener(cfg)

	if _, err := s.SpawnMeta(listener, gen.MetaOptions{}); err != nil {
		return err
	}

	s.Log().Debug("mqtt.server started (%s)", s.Name())

	return nil
}

func (s *server) HandleMessage(_ gen.PID, msg any) error {
	s.Log().Debug("mqtt.server receive message: %#v", msg)
	return inter.Send(s, msg, "message", "mqtt_producer")
}
