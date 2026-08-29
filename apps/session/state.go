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
	st.DefField("device_id", 's', "device id"),
	st.DefField("devices", 'a', "devices"),
	st.DefField("ready", 's', "ready"),
	st.DefField("error", 's', "error"),
	st.DefField("has_error", 's', "error"),
	st.DefFieldUnit("heat_loss", 'i', "heat loss", "%"),
	st.DefFieldUnit("max_power", 'i', "max power", "W"),
	st.DefFieldUnit("max_collect", 'i', "max collect", "g/h"),
	st.DefFieldUnit("net_power", 'i', "net power", "W"),
	st.DefFieldUnit("min_stability_period", 'i', "min stability period", "m"),
	st.DefField("min_reflux_ratio", 'f', "min reflux ratio"),
	st.DefField("max_stable_temp_diff", 'i', "max stable temp diff"),
	st.DefFieldUnit("body_collected_value", 'i', "body collected", "mg"),
	st.DefFieldUnit("head_collected_value", 'i', "head collected", "mg"),
	st.DefFieldUnit("recyc_collected_value", 'i', "recyc collected", "mg"),
	st.DefFieldUnit("waste_collected_value", 'i', "waste collected", "mg"),
	st.DefFieldUnit("body_collect_duration", 'i', "body collect duration", "ms"),
	st.DefFieldUnit("head_collect_duration", 'i', "head collect duration", "ms"),
	st.DefFieldUnit("recyc_collect_duration", 'i', "recyc collect duration", "ms"),
	st.DefField("collect_last_type", 's', "collect last type"),
	st.DefFieldUnit("collection_speed", 'i', "collection speed", "g/h"),
	st.DefField("reflux_ratio", 'f', "reflux ratio"),
	st.DefFieldUnit("mid_stable_temp", 'i', "middle stable temp", "d°C"),
	st.DefFieldUnit("mid_stable_first", 'i', "middle stable first temp", "d°C"),
	st.DefFieldUnit("mid_stable_diff", 'i', "middle stable first temp diff", "d°C"),
	st.DefFieldUnit("mid_stable_from", 'i', "middle stable from", "°C"),
	st.DefField("is_mid_stable", 's', "is middle temp stable"),
}

func (s *state) Init(_ ...any) error {
	s.InitState("[(session.state)]", stateFields)
	return nil
}

func (s *state) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	if kv, ok := req.(map[string]string); ok {
		s.Log().Info("session.state.HandleCall.req: %#v", req)
		return "ok", s.BulkChangeFromMap(kv, "session_repo", "session_producer")
	} else {
		return s.HandleDataRequest(req)
	}
	// return s.HandleDataRequest(req)
}

func (s *state) HandleMessage(_ gen.PID, req any) error {
	return s.HandleChangeRequests(req, "session_producer")
}

func (s *state) Terminate(reason error) {
	s.Log().Debug("session.state terminated: %s", reason)
}
