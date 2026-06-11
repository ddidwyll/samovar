package session

import (
	"samovar/lib/change"
	"samovar/lib/inter"
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

var stateFields = st.Fields{
	st.FieldParams{"device_id", 's', "device id", ""},
	st.FieldParams{"devices", 'a', "devices", ""},
	st.FieldParams{"started", 's', "started", ""},
	st.FieldParams{"error", 's', "error", ""},
	st.FieldParams{"heat_loss", 'i', "heat loss", "%"},
	st.FieldParams{"max_power", 'i', "max power", "W"},
	st.FieldParams{"max_collect", 'i', "max collect", "ml/h"},
}

func (s *state) Init(_ ...any) error {
	inter.RegisterActor(s, "[(session.state)]")
	s.data = st.New(stateFields)

	s.Log().Debug("session.state started (%s)", s.Name())
	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return s.data.HandleReq(req)
}

func (s *state) HandleMessage(_ gen.PID, msg any) error {
	s.Log().Info("session.state received message: %#v", msg)

	switch v := msg.(type) {
	case []change.Request:
		return s.updateState(v)
	case change.Request:
		return s.updateState([]change.Request{v})
	default:
		err := fmt.Sprintf("session.state unexpected message: %#v", msg)
		return errors.New(err)
	}

	return nil
}

func (s *state) updateState(reqs []change.Request) error {
	s.Log().Debug("session.state change requests: %#v", reqs)

	if reports, err := s.data.BulkChange(reqs); err != nil {
		return err
	} else {
		for _, report := range reports {
			report = report.AddFrom("session_state")
			err = inter.Send(s, report, "report", "session_producer")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("session.state terminated: %s", reason)
}
