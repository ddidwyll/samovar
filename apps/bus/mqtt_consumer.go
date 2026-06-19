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
		"mqtt_calc",
	)

	mc.AddReportRoute(
		"device_raw_state_changed",
		"mqtt_calc",
	)

	return nil
}

func (mc *mqttConsumer) HandleEvent(event gen.MessageEvent) error {
	mc.Log().Debug("bus.mqttConsumer.HandleEvent.event: %v", event)
	return mc.HandleChangeReports(event)
}
