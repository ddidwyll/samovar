package device

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type rawState struct{ st.StateActor }

var rawStateFields = st.Fields{
	st.DefFieldUnit("term_d", 'f', "t_top", "°C"),
	st.DefFieldUnit("term_c", 'f', "t_mid", "°C"),
	st.DefFieldUnit("term_k", 'f', "t_btm", "°C"),
	st.DefFieldUnit("power", 'i', "power_fact", "W"),
	st.DefFieldUnit("power_m", 'i', "power_plan", "W"),
	st.DefFieldUnit("press_a", 'f', "press_atm", "mm"),
	st.DefField("flag_otb", 's', "collect_mode"),
	st.DefField("work", 's', "work_mode"),
	st.DefFieldUnit("otbor", 'i', "collect_fact", "%"),
	st.DefFieldUnit("otbor_g_1", 'i', "cllct_head", "%"),
	st.DefFieldUnit("otbor_g_2", 'i', "cllct_ahead", "%"),
	st.DefFieldUnit("otbor_t", 'i', "cllct_body", "%"),
	st.DefFieldUnit("delta_t", 'i', "delta_body", "°C"),
	st.DefFieldUnit("time_stop", 'i', "max_stop", "s"),
	st.DefFieldUnit("otbor_minus", 'i', "decrement", "%"),
	st.DefFieldUnit("min_otb", 'i', "period_full", "m"),
	st.DefFieldUnit("sek_otb", 'i', "period_cllct", "s"),
	st.DefFieldUnit("term_d_m", 'f', "t_top_max", "°C"),
	st.DefFieldUnit("term_c_max", 'f', "t_mid_max", "°C"),
	st.DefFieldUnit("term_c_min", 'f', "t_mid_min", "°C"),
	st.DefFieldUnit("term_k_max", 'f', "t_btm_max", "°C"),
	st.DefFieldUnit("term_nasos", 'f', "water_on_t", "°C"),
	st.DefFieldUnit("term_k_m", 'f', "full_pwr_t", "°C"),
	st.DefField("kontaktor", 's', "kontaktor"),
	st.DefField("num_error", 's', "err_number"),
	st.DefFieldUnit("count_vent", 's', "count_vent", "?"),
	st.DefFieldUnit("term_vent", 'f', "term_vent", "?"),
	st.DefFieldUnit("term_v", 's', "term_v", "?"),

	// st.DefField("last_tx", 's', "last tx"),
	// st.DefField("last_rx", 's', "last rx"),
	st.DefField("last_ping", 's', "last ping"),
	st.DefField("tick", 'i', "tick"),

	st.DefField("collect_synced", 's', "is collect synced"),

	st.DefField("power_m_new", 's', "_"),
	st.DefField("otbor_new", 'i', "_"),
	st.DefField("otbor_t_new", 'i', "_"),
	st.DefField("otbor_g_1_new", 'i', "_"),
	st.DefField("otbor_g_2_new", 'i', "_"),
}

func newRawState() gen.ProcessBehavior {
	return &rawState{}
}

func (rs *rawState) Init(_ ...any) error {
	rs.InitState("[(device.raw_state)]", rawStateFields)
	return nil
}

func (rs *rawState) HandleMessage(_ gen.PID, req any) error {
	return rs.HandleChangeRequests(req, "device_producer")
}

func (rs *rawState) HandleCall(_ gen.PID, _ gen.Ref, req any) (any, error) {
	if kv, ok := req.(map[string]string); ok {
		rs.Log().Debug("device.rawState.HandleCall.req: %v", req)
		return "ok", rs.BulkChangeFromMap(kv, "mqtt_client", "device_producer")
	} else {
		return rs.HandleDataRequest(req)
	}
}
