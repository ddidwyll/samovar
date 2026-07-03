package session

import (
	clc "samovar/lib/calc"
	"samovar/lib/val"

	"fmt"
)

func calcStableMidTemp(r clc.Report, fetch clc.FetchFn, apply clc.ApplyFn) {
	currentTime := r.NewTimestamp
	currentTemp := fetch("device_state.t_mid")
	minPeriod := fetch("session_state.min_stability_period")
	if !minPeriod.IsInt() || !currentTemp.IsFlt() {
		apply("session_state.is_mid_stable", "false")
		return
	}
	stableTemp := fetch("session_state.mid_stable_temp")
	stableFrom := fetch("session_state.mid_stable_from")
	if !stableTemp.IsFlt() || !stableFrom.IsInt() || !stableTemp.Eq(currentTemp) {
		apply("session_state.is_mid_stable", "false")
		apply("session_state.mid_stable_temp", currentTemp)
		apply("session_state.mid_stable_from", currentTime)
		return
	}
	stableDuration := currentTime - stableFrom.ToInt()
	if stableDuration >= minPeriod.ToInt()*60*1000 {
		apply("session_state.is_mid_stable", "true")
	}
}

func recordCollectedValue(r clc.Report, fetch clc.FetchFn, apply clc.ApplyFn) {
	oldValue := r.OldValue
	newValue := r.NewValue
	if !oldValue.IsInt() || !newValue.IsInt() {
		return
	}

	if oldValue.EqStr("0") {
		currentType := fetch("device_state.collect_type")
		if !currentType.EqStr("OFF") {
			apply("session_state.collect_last_type", currentType)
		}
		return
	}

	currentType := fetch("session_state.collect_last_type")
	if currentType.IsNil() {
		return
	}

	fullCollectGH := fetch("session_state.max_collect")
	if !fullCollectGH.IsInt() {
		return
	}
	fullCollectMgMs := fullCollectGH.ToFlt() / 3600.0
	durationMs := r.NewTimestamp - r.OldTimestamp
	collectedMg := int64(fullCollectMgMs * oldValue.ToFlt() / 100.0 * float64(durationMs))

	// fmt.Printf("### recordCollectedValue.durationS: %fs\n", durationMs/1000.0)
	// fmt.Printf("### recordCollectedValue.collectedG: %dg (%v)\n", collectedMg/1000.0, currentType)

	var collectAcc string
	var durationAcc string
	switch currentType.String() {
	case "BODY":
		collectAcc = "session_state.body_collected_value"
		durationAcc = "session_state.body_collect_duration"
	case "HEAD":
		collectAcc = "session_state.head_collected_value"
		durationAcc = "session_state.head_collect_duration"
	case "RECYC":
		collectAcc = "session_state.recyc_collected_value"
		durationAcc = "session_state.recyc_collect_duration"
	}

	currentCollected := fetch(collectAcc)
	if currentCollected.IsNil() {
		currentCollected = val.IntAsInt(0)
	}
	currentDuration := fetch(durationAcc)
	if currentDuration.IsNil() {
		currentDuration = val.IntAsInt(0)
	}

	apply(collectAcc, currentCollected.ToInt()+collectedMg)
	fmt.Printf("recordCollectedValue.duration[%s]: %d + %d = %d\n", durationAcc, currentDuration.ToInt(), durationMs, currentDuration.ToInt()+durationMs)
	apply(durationAcc, currentDuration.ToInt()+durationMs)
}

func calcNetPower(args clc.Args, apply clc.ApplyFn) {
	loss := args.MustGet("session_state.heat_loss")
	power := args.MustGet("device_state.power")
	if !loss.IsInt() || !power.IsInt() {
		return
	}
	_ = power.ToInt() * loss.ToInt() / 100
	// fmt.Printf("calcNetPower.diff: %d", diff)
	apply("session_state.net_power", 1700)
}

func calcRefluxRatio(args clc.Args, apply clc.ApplyFn) {
	power := args.MustGet("session_state.net_power")
	collect := args.MustGet("session_state.collection_speed")
	if !power.IsInt() || !collect.IsInt() {
		return
	}
	if collect.EqStr("0") {
		apply("session_state.reflux_ratio", 999)
		return
	}
	speedVGH := power.ToFlt() * 3600.0 / 925.0
	returnVGH := speedVGH - collect.ToFlt()
	ratio := returnVGH / collect.ToFlt()
	// fmt.Printf("calcRefluxRatio.ratio: %f / %f = %f", returnVGH, collect.ToFlt(), ratio)
	apply("session_state.reflux_ratio", ratio)
}

func calcCollecionSpeed(r clc.Report, fetch clc.FetchFn, apply clc.ApplyFn) {
	fullCollectGH := fetch("session_state.max_collect")
	if !fullCollectGH.IsInt() {
		return
	}
	collectValue := fetch("device_state.collect_value")
	if !collectValue.IsInt() || collectValue.EqStr("0") {
		apply("session_state.collection_speed", 0)
		return
	}
	collectSpeedMgH := fullCollectGH.ToInt() * collectValue.ToInt() / 100
	apply("session_state.collection_speed", collectSpeedMgH)
}
