package device

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type rawState struct {
	act.Actor
}

func newRawState() gen.ProcessBehavior {
	return &rawState{}
}

func (rs *rawState) Init(_ ...any) error {
	rs.Log().Debug("device.RawState started (%s)", rs.Name())
	return nil
}

func (rs *rawState) HandleMessage(_ gen.PID, msg any) error {
	rs.Log().Debug("device.RawState receive message: %#v", msg)
	return nil
}
