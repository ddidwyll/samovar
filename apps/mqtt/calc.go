package mqtt

import (
	"samovar/common/models"
	clc "samovar/lib/calc"
	"samovar/lib/inter"
	"samovar/lib/val"

	"ergo.services/ergo/gen"
)

type calc struct{ clc.CalcActor }

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{mqtt.calc}")

	c.SetApplyFn(c.publish)

	c.Watch(
		mqttPublish,
		"device_desired_state.power",
		"device_raw_state.power_m",
	)

	return nil
}

func (c *calc) publish(key string, value any) {
	v, ok := value.(val.Val)
	if !ok {
		panic("Unexpected mqtt calc result")
	}
	msg := models.NewMqttMessage([]byte(key), []byte(v.String()))
	if err := inter.Send(c, msg, "message", "mqtt_client"); err != nil {
		panic(err)
	}
}

func mqttPublish(args clc.Args, publish clc.ApplyFn) {
	power := args["device_desired_state.power"]
	if power.IsInt() {
		publish("power_m_new", power)
	}
}
