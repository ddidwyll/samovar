package session

import (
	clc "samovar/lib/calc"

	"fmt"
)

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

	fullCollectGH := fetch("session_state.max_collect")
	if !fullCollectGH.IsInt() {
  	return
	}
	fullCollectMgMs := fullCollectGH.ToFlt() / 3600.0
	passedMs := float64(r.NewTimestamp - r.OldTimestamp)
	collectedMg := int64(fullCollectMgMs * oldValue.ToFlt() / 100.0 * passedMs)

	fmt.Printf("### recordCollectedValue.passedS: %fs\n", passedMs / 1000.0)
	fmt.Printf("### recordCollectedValue.collectedG: %dg (%v)\n", collectedMg / 1000.0, currentType)
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
