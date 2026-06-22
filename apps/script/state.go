package script

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type state struct{ st.StateActor }

func newState() gen.ProcessBehavior {
	return &state{}
}

var stateFields = st.Fields{
	st.FieldParams{"scripts", 'a', "scripts", ""},
	st.FieldParams{"script_mode", 's', "script mode", ""},
	st.FieldParams{"collect_type", 's', "collect type", ""},
	st.FieldParams{"collect_value", 'i', "collect value", ""},
}

func (s *state) Init(_ ...any) error {
	s.InitState("[(script.state)]", stateFields)
	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return s.HandleDataRequest(req)
}

func (s *state) HandleMessage(_ gen.PID, msg any) error {
	return s.HandleChangeRequests(msg, "script_producer")
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("script.state.Terminate.reason: %s", reason)
}
