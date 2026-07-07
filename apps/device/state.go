package device

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type state struct{ st.StateActor }

var stateFields = st.Fields{
	st.DefFieldUnit("t_top", 'i', "t top", "d°C"),
	st.DefFieldUnit("t_mid", 'i', "t middle", "d°C"),
	st.DefFieldUnit("t_btm", 'i', "t bottom", "d°C"),
	st.DefFieldUnit("power", 'i', "power", "W"),
	st.DefFieldUnit("power_diff", 'i', "power diff", "%"),
	st.DefFieldUnit("press", 'f', "press", "mm"),
	st.DefFieldUnit("collect", 's', "collect", "%"),
	st.DefFieldUnit("collect_value", 'f', "collect value", "%"),
	st.DefField("collect_type", 's', "collect type"),
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
