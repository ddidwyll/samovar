package mqtt

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"errors"
)

type Server struct{ act.Actor }

func newServer() gen.ProcessBehavior { return &Server{} }

func (s *Server) Init(args ...any) error {
	cfg, ok := args[0].(*config)

	if !ok {
		return errors.New("invalid mqtt client config")
	}

	listener := newListener(cfg)

	if _, err := s.SpawnMeta(listener, gen.MetaOptions{}); err != nil {
		return err
	}

	s.Log().Info("mqtt.Server started (%s)", s.Name())

	return nil
}

func (s *Server) HandleMessage(_ gen.PID, message any) error {
	s.Log().Info("mqtt.Server receive message: %#v", message)
	return nil
}
