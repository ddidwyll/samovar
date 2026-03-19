package mqtt

import "ergo.services/ergo/gen"

type Mqtt struct{}

func CreateMqttApp() gen.ApplicationBehavior {
	return &Mqtt{}
}

func (app *Mqtt) Load(node gen.Node, _ ...any) (gen.ApplicationSpec, error) {
	spec := gen.ApplicationSpec{
		Name:        "mqtt",
		Description: "MQTT",
		Mode:        gen.ApplicationModeTransient,
		Group: []gen.ApplicationMemberSpec{
			{
				Name:    "sup",
				Factory: newSup,
			},
		},
	}

	return spec, nil
}

func (app *Mqtt) Start(_ gen.ApplicationMode) {}
func (app *Mqtt) Terminate(_reason error)     {}
