package device

import (
	"samovar/lib/change"
	"samovar/lib/i"
	"samovar/lib/state"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"errors"
	"fmt"
)

type rawState struct {
	act.Actor
	data *state.State
}

func newRawState() gen.ProcessBehavior {
	return &rawState{}
}

func (rs *rawState) Init(_ ...any) error {
	rs.data = state.New(state.Fields{
		state.FieldParams{"term_d", 'f', "t_top", "°C"},
		state.FieldParams{"term_c", 'f', "t_mid", "°C"},
		state.FieldParams{"term_k", 'f', "t_btm", "°C"},
		state.FieldParams{"power", 'i', "power_fact", "W"},
		state.FieldParams{"power_m", 'i', "power_plan", "W"},
		state.FieldParams{"press_a", 'f', "press_atm", "mm"},
		state.FieldParams{"flag_otb", 's', "collect_mode", ""},
		state.FieldParams{"work", 's', "work_mode", ""},
		state.FieldParams{"otbor", 'i', "collect_fact", "%"},
		state.FieldParams{"otbor_g_1", 'i', "cllct_head", "%"},
		state.FieldParams{"otbor_g_2", 'i', "cllct_ahead", "%"},
		state.FieldParams{"otbor_t", 'i', "cllct_body", "%"},
		state.FieldParams{"delta_t", 'i', "delta_body", "°C"},
		state.FieldParams{"time_stop", 'i', "max_stop", "s"},
		state.FieldParams{"otbor_minus", 'i', "decrement", "%"},
		state.FieldParams{"min_otb", 'i', "period_full", "m"},
		state.FieldParams{"sek_otb", 'i', "period_cllct", "s"},
		state.FieldParams{"term_d_m", 'f', "t_top_max", "°C"},
		state.FieldParams{"term_c_max", 'f', "t_mid_max", "°C"},
		state.FieldParams{"term_c_min", 'f', "t_mid_min", "°C"},
		state.FieldParams{"term_k_max", 'f', "t_btm_max", "°C"},
		state.FieldParams{"term_nasos", 'f', "water_on_t", "°C"},
		state.FieldParams{"term_k_m", 'f', "full_pwr_t", "°C"},
		state.FieldParams{"kontaktor", 's', "kontaktor", ""},
		state.FieldParams{"num_error", 's', "err_number", ""},
		state.FieldParams{"count_vent", 's', "count_vent", "?"},
		state.FieldParams{"term_vent", 's', "term_vent", "?"},
		state.FieldParams{"term_v", 's', "term_v", "?"},
	})

	rs.Log().Debug("device.rawState started (%s)", rs.Name())
	return nil
}

func (rs *rawState) HandleMessage(_ gen.PID, msg any) error {
	rs.Log().Debug("device.rawState receive message: %#v", msg)

	switch v := msg.(type) {
	case change.Request:
		return rs.updateState(v)
	default:
		err := fmt.Sprintf("device.rawState unexpected message: %#v", msg)
		return errors.New(err)
	}
}

func (rs *rawState) updateState(req change.Request) error {
	rs.Log().Debug("device.rawState change req: %#v", req)
	report, err := rs.data.Change(req)

	if err == nil && report.Changed {
		report = report.AddFrom("device_raw_state")
		if err = rs.Send("device_producer", report); err != nil {
			return err
		}

		if !i.N(req.Key, "press_a", "power") {
			rs.Log().Info(report.MakeLogString("device.rawState"))
		}
	}

	// if err != nil {
	//  	rs.Log().Error("device.rawState state error: %s", err)
	// }

	return err
}
