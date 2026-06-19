package bus

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type sup struct{ act.Supervisor }

func newSup() gen.ProcessBehavior { return &sup{} }

func (s *sup) Init(_ ...any) (spec act.SupervisorSpec, err error) {
	spec.Type = act.SupervisorTypeOneForOne
	spec.Restart.Strategy = act.SupervisorStrategyTransient
	spec.Restart.Intensity = 2
	spec.Restart.Period = 5

	spec.Children = []act.SupervisorChildSpec{
		{
			Name:    "mqtt_producer",
			Factory: newMqttProducer,
		},
		{
			Name:    "device_producer",
			Factory: newDeviceProducer,
		},
		{
			Name:    "client_producer",
			Factory: newClientProducer,
		},
		{
			Name:    "session_producer",
			Factory: newSessionProducer,
		},
		{
			Name:    "mqtt_consumer",
			Factory: newMqttConsumer,
		},
		{
			Name:    "device_consumer",
			Factory: newDeviceConsumer,
		},
		{
			Name:    "client_consumer",
			Factory: newClientConsumer,
		},
		{
			Name:    "session_consumer",
			Factory: newSessionConsumer,
		},
	}

	return spec, err
}
