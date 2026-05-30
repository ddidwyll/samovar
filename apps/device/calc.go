package device

import (
	cc "samovar/lib/calc"
	"samovar/lib/val"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"

	"fmt"
)

type calc struct {
	act.Actor
	calc *cc.Calc
}

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.calc = cc.NewCalc(c)

	c.calc.Watch(
		defaultCalnFn,
		"device_state.press",
		"device_raw_state.press_a",
	)

	return nil
}

func defaultCalnFn(args cc.Args) (val.Val, error) {
	fmt.Println("%v", args)
	return val.Nil{}, nil
}
