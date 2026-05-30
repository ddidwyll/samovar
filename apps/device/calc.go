package device

import (
	clc "samovar/lib/calc"
	"samovar/lib/val"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"fmt"
)

type calc struct {
	act.Actor
	config *clc.Config
}

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.config = clc.NewConfig(c)

	c.config.Watch(
		asIs,
		"device_state.t_top",
		"device_raw_state.term_d",
	)

	c.config.Watch(
		asIs,
		"device_state.t_mid",
		"device_raw_state.term_c",
	)

	c.config.Watch(
		asIs,
		"device_state.t_btm",
		"device_raw_state.term_k",
	)

	c.config.Watch(
		asIs,
		"device_state.press",
		"device_raw_state.press_a",
	)

	c.config.Watch(
		asIs,
		"device_state.power",
		"device_raw_state.power_m",
	)

	c.config.Watch(
		calcPowerDiff,
		"device_state.power_diff",
		"device_raw_state.power_m",
		"device_raw_state.power",
	)

	c.config.Watch(
		calcCollect,
		"device_state.collect",
		"device_raw_state.otbor",
		"device_raw_state.flag_otb",
		"device_raw_state.sek_otb",
		"device_raw_state.min_otb",
	)

	return nil
}

func calcPowerDiff(args clc.Args) (val.Val, error) {
	power_fact := args["device_raw_state.power"]
	power_plan := args["device_raw_state.power_m"]

	if !power_fact.IsInt() || !power_plan.IsInt() {
		return val.Nil{}, nil
	}

	plan := power_plan.ToInt()
	fact := power_fact.ToInt()
	diff := (plan - fact) * 100 / plan
	return val.IntAsInt(diff), nil
}

func calcCollect(args clc.Args) (val.Val, error) {
	f, s := 0.0, "OFF"

	otbor := args["device_raw_state.otbor"]
	if !otbor.IsInt() {
		return val.Nil{}, nil
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

func asIs(args clc.Args) (val.Val, error) {
	for _, val := range args {
		return val, nil
	}
	return val.Nil{}, nil
}

func (c *calc) HandleMessage(_ gen.PID, changeReport any) error {
	return c.config.HandleReport(changeReport)
}
