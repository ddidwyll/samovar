package device

import (
	clc "samovar/lib/calc"
	"samovar/lib/val"

	"fmt"
)

func calcPowerDiff(args clc.Args, apply clc.ApplyFn) {
	power_fact := args.MustGet("device_raw_state.power")
	power_plan := args.MustGet("device_raw_state.power_m")

	if !power_fact.IsInt() || !power_plan.IsInt() {
		return
	}

	plan := power_plan.ToInt()
	fact := power_fact.ToInt()
	diff := (plan - fact) * 100 / plan

	apply("device_state.power_diff", val.IntAsInt(diff))
}

func calcCollect(args clc.Args, apply clc.ApplyFn) {
	f, s := 0.0, "OFF"

	otbor := args.MustGet("device_raw_state.otbor")
	if !otbor.IsInt() {
		return
	}

	calcPeriod := func() float64 {
		sec := args.MustGet("device_raw_state.sek_otb")
		min := args.MustGet("device_raw_state.min_otb")

		if !sec.IsInt() || !min.IsInt() {
			return -1.0
		}

		return sec.ToFlt() * 100.0 / (min.ToFlt() * 60.0)
	}

	switch args.MustGet("device_raw_state.flag_otb").String() {
	case "Golov":
		f, s = otbor.ToFlt(), "HEAD"
	case "Podgol":
		f, s = otbor.ToFlt(), "RECYC"
	case "Gol.P":
		f, s = calcPeriod(), "SALVO"
	case "Telo":
		f, s = otbor.ToFlt(), "BODY"
	}

	collectValue := val.FltAsFlt(f)
	collectType, _ := val.StrAsStr(s)
	collectStr := fmt.Sprintf("%s_%s", collectType, collectValue)
	collect, _ := val.StrAsStr(collectStr)

	apply("device_state.collect_value", collectValue)
	apply("device_state.collect_type", collectType)
	apply("device_state.collect", collect)
}
