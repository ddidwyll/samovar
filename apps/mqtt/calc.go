package mqtt

import (
	clc "samovar/lib/calc"
	"samovar/lib/val"

	"ergo.services/ergo/gen"
)

type calc struct{ clc.CalcActor }

func newCalc() gen.ProcessBehavior {
	return &calc{}
}

func (c *calc) Init(_ ...any) error {
	c.InitCalc("{mqtt.calc}")

	c.Watch(translateToMqtt, "device_desired_state.power")

	return nil
}

func mqttPublish(topic string, value string) 

func translateToMqtt(args clc.Args, apply clc.ApplyFn) {
	power := args["device_desired_state.power"]

	if !power.IsInt() {
		return
	}
}
