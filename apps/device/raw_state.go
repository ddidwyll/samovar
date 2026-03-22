package device

import (
	"samovar/lib/change"
	"samovar/lib/state"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
	"fmt"
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
	case change.Request:
		return rs.updateState(v)
	default:
		err := fmt.Sprintf("device.rawState unexpected message: %#v", msg)
		return errors.New(err)
	}
}

func (rs *rawState) updateState(req change.Request) error {
	rs.Log().Debug("device.rawState change req: %#v", req)
	report, err := rs.state.Change(req)

	if err == nil && report.Changed {
		rs.Log().Info("device.rawState changed: %s -> %s", report.OldValue.ToStr(), report.NewValue.ToStr())
	}

	// if err != nil {
	//  	rs.Log().Error("device.rawState state error: %s", err)
	// }

	return nil
}
