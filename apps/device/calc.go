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
	calc *clc.Config
}

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.calc = clc.NewCalc(c)

	c.calc.Watch(
		defaultCalnFn,
		"device_state.press",
		"device_raw_state.press_a",
	)

	return nil
}

func defaultCalnFn(args clc.Args) (val.Val, error) {
	fmt.Println("%v", args)
	return val.Nil{}, nil
}

func (c *calc) HandleCall(_ gen.PID, _ gen.Ref, changeReport any) (any, error) {
	return nil, c.calc.HandleReport(changeReport)
}
