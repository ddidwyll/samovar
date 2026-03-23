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
		state.FieldParams{"term_d", 'f', "t top", "°C"},
		state.FieldParams{"term_c", 'f', "t middle", "°C"},
		state.FieldParams{"term_k", 'f', "t bottom", "°C"},
		state.FieldParams{"power", 'i', "power", "W"},
		state.FieldParams{"press_a", 'f', "atm press", "mm"},
		state.FieldParams{"flag_otb", 's', "collect mode", ""},
		state.FieldParams{"term_d_m", 'f', "t top max", "°C"},
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
		field, _ := rs.state.Fetch(req.Key)

		if report.IsNew {
			rs.Log().Info("device.rawState [%s\t]:\t%s", field, report.NewValue)
		} else {
			rs.Log().Info("device.rawState [%s\t]:\t%s -> %s\t\t%s\t%s", field, report.OldValue, report.NewValue, report.FormatTS(), report.Reason)
		}
	}

	// if err != nil {
	//  	rs.Log().Error("device.rawState state error: %s", err)
	// }

	return nil
}
