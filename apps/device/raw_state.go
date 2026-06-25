package device

import (
	st "samovar/lib/state"

	"ergo.services/ergo/gen"
)

type rawState struct{ st.StateActor }

var rawStateFields = st.Fields{
	st.FieldParams{"term_d", 'f', "t_top", "°C"},
	st.FieldParams{"term_c", 'f', "t_mid", "°C"},
	st.FieldParams{"term_k", 'f', "t_btm", "°C"},
	st.FieldParams{"power", 'i', "power_fact", "W"},
	st.FieldParams{"power_m", 'i', "power_plan", "W"},
	st.FieldParams{"press_a", 'f', "press_atm", "mm"},
	st.FieldParams{"flag_otb", 's', "collect_mode", ""},
	st.FieldParams{"work", 's', "work_mode", ""},
	st.FieldParams{"otbor", 'i', "collect_fact", "%"},
	st.FieldParams{"otbor_g_1", 'i', "cllct_head", "%"},
	st.FieldParams{"otbor_g_2", 'i', "cllct_ahead", "%"},
	st.FieldParams{"otbor_t", 'i', "cllct_body", "%"},
	st.FieldParams{"delta_t", 'i', "delta_body", "°C"},
	st.FieldParams{"time_stop", 'i', "max_stop", "s"},
	st.FieldParams{"otbor_minus", 'i', "decrement", "%"},
	st.FieldParams{"min_otb", 'i', "period_full", "m"},
	st.FieldParams{"sek_otb", 'i', "period_cllct", "s"},
	st.FieldParams{"term_d_m", 'f', "t_top_max", "°C"},
	st.FieldParams{"term_c_max", 'f', "t_mid_max", "°C"},
	st.FieldParams{"term_c_min", 'f', "t_mid_min", "°C"},
	st.FieldParams{"term_k_max", 'f', "t_btm_max", "°C"},
	st.FieldParams{"term_nasos", 'f', "water_on_t", "°C"},
	st.FieldParams{"term_k_m", 'f', "full_pwr_t", "°C"},
	st.FieldParams{"kontaktor", 's', "kontaktor", ""},
	st.FieldParams{"num_error", 's', "err_number", ""},
	st.FieldParams{"count_vent", 's', "count_vent", "?"},
	st.FieldParams{"term_vent", 'f', "term_vent", "?"},
	st.FieldParams{"term_v", 's', "term_v", "?"},

	st.FieldParams{"last_tx", 's', "last tx", ""},

	st.FieldParams{"power_m_new", 's', "_", ""},
	st.FieldParams{"otbor_new", 'i', "_", ""},
	st.FieldParams{"otbor_t_new", 'i', "_", ""},
	st.FieldParams{"otbor_g_1_new", 'i', "_", ""},
	st.FieldParams{"otbor_g_2_new", 'i', "_", ""},
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
	return rs.HandleDataRequest(req)
}
