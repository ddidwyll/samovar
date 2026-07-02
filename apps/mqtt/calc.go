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

	c.WatchFields(
		changePower,
		"device_desired_state.power",
		"device_raw_state.power_m",
	)

	c.WatchFields(
		changeCollect,
		"device_raw_state.last_ping",
		"device_raw_state.otbor",
		"device_state.collect_type",
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
	if clc.IsFid(key) {
		if err := c.SendRequest(key, value); err != nil {
			panic(err)
		}
	} else {
		msg := models.NewMqttMessage([]byte(key), []byte(strVal))
		if err := inter.Send(c, msg, "message", "mqtt_publisher"); err != nil {
			panic(err)
		}
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
	newType := args.MustGet("device_desired_state.collect_type").String()
	newValue := args.MustGet("device_desired_state.collect_value")

	if newType == "" || !newValue.IsInt() {
		return
	}

	currentType := args.MustGet("device_state.collect_type").String()
	currentValue := args.MustGet("device_raw_state.otbor")

	if currentType == "" || !currentValue.IsInt() {
		return
	}

	if newType != currentType && newType != "OFF" && currentType != "OFF" {
		newType = "TEMP_OFF"
	}

	switch newType {
	case "BODY":
		if currentType != "BODY" {
			change("device_raw_state.collect_synced", "false")
			change("otbor_t_new", newValue)
			change("work", 8)
		} else if !currentValue.Eq(newValue) {
			change("device_raw_state.collect_synced", "false")
			change("otbor_new", newValue)
		} else {
			change("device_raw_state.collect_synced", "true")
		}
	case "HEAD":
		if currentType != "HEAD" {
			change("device_raw_state.collect_synced", "false")
			change("otbor_g_1_new", newValue)
			change("work", 9)
		} else if !currentValue.Eq(newValue) {
			change("device_raw_state.collect_synced", "false")
			change("otbor_new", newValue)
		} else {
			change("device_raw_state.collect_synced", "true")
		}
	case "RECYC":
		if currentType != "RECYC" {
			change("device_raw_state.collect_synced", "false")
			change("otbor_g_2_new", newValue)
			change("work", 10)
		} else if !currentValue.Eq(newValue) {
			change("device_raw_state.collect_synced", "false")
			change("otbor_new", newValue)
		} else {
			change("device_raw_state.collect_synced", "true")
		}
	default:
		if currentType != "OFF" {
			change("device_raw_state.collect_synced", "false")
			change("work", 6)
		} else {
			change("device_raw_state.collect_synced", "true")
		}
	}
}

func (c *calc) HandleMessage(_ gen.PID, msg any) error {
	c.Log().Debug("mqtt.calc.HandleMessage.msg: %v", msg)
	return c.HandleChangeReports(msg)
}
