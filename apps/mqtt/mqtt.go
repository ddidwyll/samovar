package mqtt

import (
	"ergo.services/ergo/gen"
)

func CreateMqtt() gen.ApplicationBehavior {
	return &Mqtt{}
}

type Mqtt struct{}

// Load invoked on loading application using method ApplicationLoad of gen.Node interface.
func (app *Mqtt) Load(node gen.Node, args ...any) (gen.ApplicationSpec, error) {
	return gen.ApplicationSpec{
		Name:        "mqtt",
		Description: "description of this application",
		Mode:        gen.ApplicationModeTransient,
		Group: []gen.ApplicationMemberSpec{
			{
				Name:    "mqttsup",
				Factory: factory_MqttSup,
			},
		},
	}, nil
}

// Start invoked once the application started
func (app *Mqtt) Start(mode gen.ApplicationMode) {}

// Terminate invoked once the application stopped
func (app *Mqtt) Terminate(reason error) {}
