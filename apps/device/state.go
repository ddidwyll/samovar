package device

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type state struct{ st.StateActor }

var stateFields = st.Fields{
	st.FieldParams{"t_top", 'i', "t top", "d°C"},
	st.FieldParams{"t_mid", 'i', "t middle", "d°C"},
	st.FieldParams{"t_btm", 'i', "t bottom", "d°C"},
	st.FieldParams{"power", 'i', "power", "W"},
	st.FieldParams{"power_diff", 'i', "power diff", "%"},
	st.FieldParams{"collect", 's', "collect", "%"},
	st.FieldParams{"collect_value", 'f', "collect value", "%"},
	st.FieldParams{"collect_type", 's', "collect type", ""},
	st.FieldParams{"press", 'f', "press", "mm"},
}

func newState() gen.ProcessBehavior {
	return &state{}
}

func (s *state) Init(_ ...any) error {
	s.InitState("[(device.state)]", stateFields)
	return nil
}

func (s *state) HandleMessage(_ gen.PID, req any) error {
	return s.HandleChangeRequests(req, "device_producer")
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	return s.HandleDataRequest(req)
}
