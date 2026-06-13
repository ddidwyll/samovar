package session

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type state struct{ st.StateActor }

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
	s.InitState("[(session.state)]", stateFields)
	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return s.HandleDataRequest(req)
}

func (s *state) HandleMessage(_ gen.PID, req any) error {
	return s.HandleChangeRequests(req, "session_producer")
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("session.state terminated: %s", reason)
}
