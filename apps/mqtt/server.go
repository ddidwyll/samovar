package mqtt

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Server struct {
	act.Actor
}

func newServer() gen.ProcessBehavior {
	return &Server{}
}

func (s *Server) Init(_ ...any) error {
	listener := newListener()

	if _, err := s.SpawnMeta(listener, gen.MetaOptions{}); err != nil {
		s.Log().Error("mqtt.Server error: %s", err)
		return err
	}

	s.Log().Info("mqtt.Server started (%s)", s.Name())

	return nil
}

func (s *Server) HandleMessage(_ gen.PID, message any) error {
	s.Log().Info("mqtt.Server receive message: %#v", message)
	return nil
}
