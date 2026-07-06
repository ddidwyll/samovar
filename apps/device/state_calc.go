package device

import (
	clc "samovar/lib/calc"
	"samovar/lib/val"

	"fmt"
	"time"
)

func calcDecigrad(args clc.Args, apply clc.ApplyFn) {
	mapper := map[string]string{
		"term_d": "t_top",
		"term_c": "t_mid",
		"term_k": "t_btm",
	}

	for key, val := range args {
		target := clc.Fid(key)
		targetField, exists := mapper[target.FieldKey]
		if !val.IsFlt() || !exists {
			continue
		}
		targetKey := "device_state." + targetField
		apply(targetKey, int64(val.ToFlt()*10.0))
	}
}

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

	randDecimal := float64(time.Now().Second()/6) / 100
	otborRand := otbor.ToFlt() + randDecimal

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
		f, s = otborRand, "HEAD"
	case "Podgol":
		t := "RECYC"
		if otbor.ToInt() >= 50 {
			t = "WASTE"
		}
		f, s = otborRand, t
	case "Gol.P":
		f, s = calcPeriod(), "SALVO"
	case "Telo":
		f, s = otborRand, "BODY"
	}

	collectValue := val.FltAsFlt(f)
	collectType, _ := val.StrAsStr(s)
	collectStr := fmt.Sprintf("%s_%d", collectType, collectValue.ToInt())
	collect, _ := val.StrAsStr(collectStr)

	apply("device_state.collect_value", collectValue)
	apply("device_state.collect_type", collectType)
	apply("device_state.collect", collect)
}
