package mqtt

import "ergo.services/ergo/gen"

type mqtt struct{}

func CreateMqttApp() gen.ApplicationBehavior {
	return &mqtt{}
}

func (app *mqtt) Load(_ gen.Node, _ ...any) (spec gen.ApplicationSpec, err error) {
	cfg, err := loadConfig()

	spec.Name = "mqtt_app"
	spec.Description = "MQTT"
	spec.Mode = gen.ApplicationModeTransient
	spec.Group = []gen.ApplicationMemberSpec{
		{
			Name:    "mqtt_sup",
			Factory: newSup,
			Args:    []any{cfg},
		},
	}

	return spec, err
}

func (app *mqtt) Start(_ gen.ApplicationMode) {}
func (app *mqtt) Terminate(_ error)           {}
