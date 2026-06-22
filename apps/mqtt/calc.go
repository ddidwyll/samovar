package mqtt

import (
	"samovar/common/models"
	clc "samovar/lib/calc"
	"samovar/lib/inter"
	"samovar/lib/val"

	"ergo.services/ergo/gen"

	"fmt"
)

type calc struct{ clc.CalcActor }

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{mqtt.calc}")

	c.SetApplyFn(c.publish)

	c.Watch(
		changePower,
		"device_desired_state.power",
		"device_raw_state.power_m",
	)

	c.Watch(
		changeCollect,
		"device_desired_state.collect_type",
		"device_desired_state.collect_value",
	)

	return nil
}

func (c *calc) publish(key string, value any) {
	var strVal string
	switch v := value.(type) {
	case val.Val:
		strVal = v.String()
	case string:
		strVal = v
	case int64:
		strVal = fmt.Sprintf("%d", v)
	default:
		panic("Unexpected mqtt calc value")
	}
	msg := models.NewMqttMessage([]byte(key), []byte(strVal))
	if _, err := inter.Call(c, msg, "message", "mqtt_client"); err != nil {
		panic(err)
	}
}

func changePower(args clc.Args, change clc.ApplyFn) {
	power := args.MustGet("device_desired_state.power")
	current_power := args.MustGet("device_raw_state.power_m")
	if power.IsInt() && !current_power.Eq(power) {
		change("power_m_new", power)
	}
}

func changeCollect(args clc.Args, change clc.ApplyFn) {
	ctype := args.MustGet("device_desired_state.collect_type")
	cval := args.MustGet("device_desired_state.collect_value")

	if ctype.IsNil() || !cval.IsInt() {
		return
	}

	typeS := ctype.String()
	valI := cval.ToInt()
	if valI <= 0 || valI >= 100 {
		typeS = "OFF"
	}
	if typeS == "BODY" {
		change("otbor_t", valI)
		change("work", 7)
		return
	}
	if typeS == "HEAD" {
		change("otbor_g_1", valI)
		change("work", 9)
		return
	}
	if typeS == "RECYC" {
		change("otbor_g_2", valI)
		change("work", 10)
		return
	}
	change("work", 6)
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	c.Log().Debug("mqtt.calc.HandleMessage.msg: %v", msg)
	return c.HandleChangeReports(msg)
}
