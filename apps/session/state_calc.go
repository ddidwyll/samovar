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
	fmt.Printf("### recordCollectedValue.Report: (%v) %v\n", currentType, r)
}
