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
	st.FieldParams{"ready", 's', "ready", ""},
	st.FieldParams{"error", 's', "error", ""},
	st.FieldParams{"has_error", 's', "error", ""},
	st.FieldParams{"heat_loss", 'i', "heat loss", "%"},
	st.FieldParams{"max_power", 'i', "max power", "W"},
	st.FieldParams{"max_collect", 'i', "max collect", "g/h"},
	st.FieldParams{"net_power", 'i', "net power", "W"},
	st.FieldParams{"min_stability_period", 'i', "min stability period", "m"},
	st.FieldParams{"min_reflux_ratio", 'f', "min reflux ratio", ""},
	st.FieldParams{"max_stable_temp_diff", 'i', "max stable temp diff", ""},
	st.FieldParams{"body_collected_value", 'i', "body collected", "mg"},
	st.FieldParams{"head_collected_value", 'i', "head collected", "mg"},
	st.FieldParams{"recyc_collected_value", 'i', "recyc collected", "mg"},
	st.FieldParams{"waste_collected_value", 'i', "waste collected", "mg"},
	st.FieldParams{"body_collect_duration", 'i', "body collect duration", "ms"},
	st.FieldParams{"head_collect_duration", 'i', "head collect duration", "ms"},
	st.FieldParams{"recyc_collect_duration", 'i', "recyc collect duration", "ms"},
	st.FieldParams{"collect_last_type", 's', "collect last type", ""},
	st.FieldParams{"collection_speed", 'i', "collection speed", "g/h"},
	st.FieldParams{"reflux_ratio", 'f', "reflux ratio", ""},
	st.FieldParams{"mid_stable_temp", 'i', "middle stable temp", "d°C"},
	st.FieldParams{"mid_stable_first", 'i', "middle stable first temp", "d°C"},
	st.FieldParams{"mid_stable_diff", 'i', "middle stable first temp diff", "d°C"},
	st.FieldParams{"mid_stable_from", 'i', "middle stable from", "°C"},
	st.FieldParams{"is_mid_stable", 's', "is middle temp stable", ""},
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
