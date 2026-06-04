package device

import (
	"samovar/lib/change"
	st "samovar/lib/state"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type state struct {
	act.Actor
	data *st.State
}

func newState() gen.ProcessBehavior {
	return &state{}
}

func (s *state) Init(_ ...any) error {
	s.data = st.New(st.Fields{
		st.FieldParams{"t_top", 'f', "t top", "°C"},
		st.FieldParams{"t_mid", 'f', "t middle", "°C"},
		st.FieldParams{"t_btm", 'f', "t bottom", "°C"},
		st.FieldParams{"power", 'i', "power", "W"},
		st.FieldParams{"power_diff", 'i', "power diff", "%"},
		st.FieldParams{"collect", 's', "collect", "%"},
		st.FieldParams{"press", 'f', "press", "mm"},
	})

	s.Log().Debug("device.state started (%s)", s.Name())
	return nil
}

func (s *state) HandleMessage(_ gen.PID, msg any) error {
	s.Log().Debug("device.state received message: %#v", msg)

	switch r := msg.(type) {
	case change.Request:
		s.Log().Warning("device.state received change request from %s", r.LastFrom())
	case []change.Request:
		return s.updateState(r)
	default:
		err := fmt.Sprintf("device.state unexpected message: %#v", msg)
		return errors.New(err)
	}
	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return s.data.HandleReq(req)
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("device.state terminated: %s", reason)
}

func (s *state) updateState(reqs []change.Request) error {
	s.Log().Debug("device.state change requests: %#v", reqs)

	if reports, err := s.data.BulkChange(reqs); err != nil {
		return err
	} else {
		for _, report := range reports {
			report = report.AddFrom("device_state")
			if err = s.Send("device_producer", report); err != nil {
				return err
			} else {
				s.Log().Debug(report.MakeLogString("device.state"))
			}
		}
	}

	return nil
}
