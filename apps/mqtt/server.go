package mqtt

import (
  "ergo.services/ergo/gen"
  "ergo.services/ergo/act"
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

  s.Log().Error("mqtt.Server started (%s)", s.Name())

  return nil
}
