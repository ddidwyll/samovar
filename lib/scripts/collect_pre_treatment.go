package scripts

func PerformCollectPreTreatment(args args, apply applyFn) {
	apply("collect_type", "BODY")
	apply("collect_value", 0)
	isStable := args.MustGet("session_state.is_mid_stable")
	if !isStable.EqStr("true") {
		// debug(isStable.String(), "session_state.is_mid_stable")
		return
	}
	diff := args.MustGet("session_state.mid_stable_diff")
	if diff.IsInt() && diff.ToInt() > 2 {
		apply("client_desired_state.script_mode", "idle")
	}
	bottomTemp := args.MustGet("device_state.t_btm")
	if !bottomTemp.IsInt() {
		debug(bottomTemp.String(), "device_state.t_btm")
		return
	}
	refluxRatio := args.MustGet("session_state.min_reflux_ratio")
	if !refluxRatio.IsFlt() || refluxRatio.ToInt() <= 0 {
		debug(refluxRatio.String(), "session_state.min_reflux_ratio")
		return
	}
	netPower := args.MustGet("session_state.net_power")
	if !netPower.IsInt() {
		debug(netPower.String(), "session_state.net_power")
		return
	}
	maxCollect := args.MustGet("session_state.max_collect")
	if !maxCollect.IsInt() {
		debug(maxCollect.String(), "session_state.max_collect")
		return
	}
	coef := 1.0
	switch {
	case bottomTemp.ToInt() > 90:
		coef = 2.0
	case bottomTemp.ToInt() > 96:
		coef = 3.0
	}
	collectGH := calcCollectGHByRefluxRatio(netPower.ToFlt(), refluxRatio.ToFlt()*coef)
	collectPerc := calcCollectPercByCollectGH(collectGH, maxCollect.ToInt())
	if collectPerc > 0 && collectPerc < 50 {
		apply("collect_value", collectPerc)
	} else {
		apply("client_desired_state.script_mode", "idle")
		return
	}
}
