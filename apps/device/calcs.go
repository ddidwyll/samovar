package device

import (
	clc "samovar/lib/calc"
	"samovar/lib/val"

	"fmt"
)

func calcPowerDiff(args clc.Args) (clc.Value, error) {
	power_fact := args["device_raw_state.power"]
	power_plan := args["device_raw_state.power_m"]

	if !power_fact.IsInt() || !power_plan.IsInt() {
		return clc.Skip()
	}

	plan := power_plan.ToInt()
	fact := power_fact.ToInt()
	diff := (plan - fact) * 100 / plan
	return val.IntAsInt(diff), nil
}

func calcCollect(args clc.Args) (clc.Value, error) {
	f, s := 0.0, "OFF"

	otbor := args["device_raw_state.otbor"]
	if !otbor.IsInt() {
		return clc.Skip()
	}

	calcPeriod := func() float64 {
		sec := args["device_raw_state.sek_otb"]
		min := args["device_raw_state.min_otb"]

		if !sec.IsInt() || !min.IsInt() {
			return -1.0
		}

		return sec.ToFlt() * 100.0 / (min.ToFlt() * 60.0)
	}

	switch args["device_raw_state.flag_otb"].String() {
	case "Golov":
		f, s = otbor.ToFlt(), "Hdrp"
	case "Podgol":
		f, s = otbor.ToFlt(), "subH"
	case "Gol.P":
		f, s = calcPeriod(), "Hprd"
	case "Telo":
		f, s = otbor.ToFlt(), "Body"
	}

	numStr := val.FltAsFlt(f).String()
	valStr := fmt.Sprintf("%s_%s", s, numStr)
	return val.StrAsStr(valStr)
}
