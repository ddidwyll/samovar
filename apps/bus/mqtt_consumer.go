package bus

import (
	"samovar/lib/stage"

	"ergo.services/ergo/gen"
)

type mqttConsumer struct{ stage.Consumer }

func newMqttConsumer() gen.ProcessBehavior {
	return &mqttConsumer{}
}

func (mc *mqttConsumer) Init(_ ...any) error {
	mc.InitConsumer("([bus.mqtt.consumer])")

	mc.AddReportRoute(
		"device_desired_state_changed",
		"mqtt_producer",
		"mqtt_calc",
	)

	return nil
}

func (mc *mqttConsumer) HandleEvent(event gen.MessageEvent) error {
	return mc.HandleChangeReports(event)
}
