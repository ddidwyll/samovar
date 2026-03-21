package device

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type RawState struct {
	act.Actor
}

func newRawState() gen.ProcessBehavior {
	return &RawState{}
}

func (rs *RawState) Init(_ ...any) error {
	rs.Log().Info("device.RawState started (%s)", rs.Name())
	return nil
}

func (rs *RawState) HandleMessage(_ gen.PID, message any) error {
	rs.Log().Info("device.RawState receive message: %#v", message)
	return nil
}
