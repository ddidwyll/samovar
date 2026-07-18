package scripts

func PerformCollectRecycFast(args args, apply applyFn) {
	refluxRatio := args.MustGet("session_state.min_reflux_ratio")
	maxCollect := args.MustGet("session_state.max_collect")
	netPower := args.MustGet("session_state.net_power")
	if !refluxRatio.IsFlt() || !maxCollect.IsInt() || !netPower.IsInt() {
		return
	}
	collectGH := calcCollectGHByRefluxRatio(netPower.ToFlt(), refluxRatio.ToFlt())
	collectPerc := calcCollectPercByCollectGH(collectGH, maxCollect.ToInt())
	if collectPerc > 0 && collectPerc < 50 {
		apply("collect_type", "RECYC")
		apply("collect_value", collectPerc)
	} else {
		apply("client_desired_state.script_mode", "idle")
		return
	}
	diff := args.MustGet("session_state.mid_stable_diff")
	if diff.IsInt() && diff.ToInt() > 2 {
		apply("client_desired_state.script_mode", "idle")
	}
}

func PerformCollectRecycSlow(args args, apply applyFn) {
	refluxRatio := args.MustGet("session_state.min_reflux_ratio")
	maxCollect := args.MustGet("session_state.max_collect")
	netPower := args.MustGet("session_state.net_power")
	if !refluxRatio.IsFlt() || !maxCollect.IsInt() || !netPower.IsInt() {
		return
	}
	collectGH := calcCollectGHByRefluxRatio(netPower.ToFlt(), refluxRatio.ToFlt()*2.0)
	collectPerc := calcCollectPercByCollectGH(collectGH, maxCollect.ToInt())
	if collectPerc > 0 && collectPerc < 50 {
		apply("collect_type", "RECYC")
		apply("collect_value", collectPerc)
	} else {
		apply("client_desired_state.script_mode", "idle")
		return
	}
	diff := args.MustGet("session_state.mid_stable_diff")
	if diff.IsInt() && diff.ToInt() > 2 {
		apply("client_desired_state.script_mode", "idle")
	}
}
