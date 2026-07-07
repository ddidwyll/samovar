package scripts

func PerformCollectBody(args args, apply applyFn) {
	apply("collect_type", "BODY")
	apply("collect_value", 0)
	isStable := args.MustGet("session_state.is_mid_stable")
	if !isStable.EqStr("true") {
		debug(isStable.String(), "session_state.is_mid_stable")
		return
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
	stableTemp := args.MustGet("session_state.mid_stable_temp")
	if !stableTemp.IsInt() {
		debug(stableTemp.String(), "session_state.mid_stable_temp")
		return
	}
	maxStableDiff := args.MustGet("session_state.max_stable_temp_diff")
	if !maxStableDiff.IsInt() {
		debug(maxStableDiff.String(), "session_state.max_stable_temp_diff")
		return
	}
	stableDiff := args.MustGet("session_state.mid_stable_diff")
	if !stableDiff.IsInt() {
		debug(stableDiff.String(), "session_state.mid_stable_diff")
		return
	}
	if stableDiff.ToInt() > maxStableDiff.ToInt() {
		debug(stableDiff.ToInt()-maxStableDiff.ToInt(), "stableDiff.ToInt() - maxStableDiff.ToInt()")
		apply("client_desired_state.script_mode", "idle")
		return
	}
	btmTempFlt := bottomTemp.ToFlt()
	if btmTempFlt < 800.0 {
		btmTempFlt = 800.0
		btmTempFlt = btmTempFlt / 10.0
	}
	collectGH := 0.174 * (100.0 - btmTempFlt) * netPower.ToFlt() / (refluxRatio.ToFlt() + 1.0)
	debug(collectGH, "collectGH")
	collectPerc := calcCollectPercByCollectGH(int64(collectGH), maxCollect.ToInt())
	debug(collectPerc, "collectPerc")
	apply("collect_type", "BODY")
	apply("collect_value", collectPerc)
}
