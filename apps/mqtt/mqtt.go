package mqtt

import "ergo.services/ergo/gen"

type Mqtt struct{}

func CreateMqttApp() gen.ApplicationBehavior {
	return &Mqtt{}
}

func (app *Mqtt) Load(node gen.Node, _ ...any) (spec gen.ApplicationSpec, err error) {
	cfg, err := loadConfig()
	if err != nil {
		return spec, err
	}

	spec.Name = "mqtt"
	spec.Description = "MQTT"
	spec.Mode = gen.ApplicationModeTransient
	spec.Group = []gen.ApplicationMemberSpec{
		{
			Name:    "sup",
			Factory: newSup,
			Args:    []any{cfg},
		},
	}

	return spec, nil
}

func (app *Mqtt) Start(_ gen.ApplicationMode) {}
func (app *Mqtt) Terminate(_ error)           {}
