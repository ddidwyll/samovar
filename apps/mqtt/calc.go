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
		"device_raw_state.last_tx",
		"device_state.collect_type",
		"device_state.collect_value",
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
	case int:
		strVal = fmt.Sprintf("%d", v)
	case int64:
		strVal = fmt.Sprintf("%d", v)
	default:
		err := fmt.Sprintf("Unexpected mqtt calc value (%T)", value)
		panic(err)
	}
	msg := models.NewMqttMessage([]byte(key), []byte(strVal))
	if err := inter.Send(c, msg, "message", "mqtt_publisher"); err != nil {
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
	newType := args.MustGet("device_desired_state.collect_type")
	newValue := args.MustGet("device_desired_state.collect_value")

	if newType.IsNil() || !newValue.IsInt() {
		return
	}

	currentType := args.MustGet("device_state.collect_type").String()
	currentValue := args.MustGet("device_state.collect_value")

	typeS := newType.String()
	valI := newValue.ToInt()
	if valI <= 0 || valI >= 100 {
		typeS = "OFF"
	}
	if typeS == "BODY" {
		if !currentValue.Eq(newValue) {
			change("otbor_t_new", valI)
			change("otbor_new", valI)
		}
		if currentType != "BODY" {
			change("work", 8)
		}
		return
	}
	if typeS == "HEAD" {
		if !currentValue.Eq(newValue) {
			change("otbor_g_1_new", valI)
			change("otbor_new", valI)
		}
		if currentType != "HEAD" {
			change("work", 9)
		}
		return
	}
	if typeS == "RECYC" {
		if !currentValue.Eq(newValue) {
			change("otbor_g_2_new", valI)
			change("otbor_new", valI)
		}
		if currentType != "RECYC" {
			change("work", 10)
		}
		return
	}
	if currentType != "OFF" {
		change("work", 6)
	}
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	c.Log().Debug("mqtt.calc.HandleMessage.msg: %v", msg)
	return c.HandleChangeReports(msg)
}
