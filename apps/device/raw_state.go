package device

import (
	"samovar/common/models"
	"samovar/lib/state"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type rawState struct {
	act.Actor
	state *state.State
}

func newRawState() gen.ProcessBehavior {
	return &rawState{}
}

func (rs *rawState) Init(_ ...any) error {
	rs.state = state.New(state.Fields{
		state.FieldParams{"term_d", 'f', "TOP temp"},
		state.FieldParams{"term_c", 'f', "MIDDLE temp"},
		state.FieldParams{"term_k", 'f', "BOTTOM temp"},
	})

	rs.Log().Debug("device.rawState started (%s)", rs.Name())
	return nil
}

func (rs *rawState) HandleMessage(_ gen.PID, msg any) error {
	rs.Log().Debug("device.rawState receive message: %#v", msg)

	switch v := msg.(type) {
	case models.ChangeRequest:
	}
	return nil
}

func (rs *rawState) updateState(key state.Key, val any) error {

}
