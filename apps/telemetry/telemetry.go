package telemetry

import "ergo.services/ergo/gen"

type telemetry struct{}

func CreateApp() gen.ApplicationBehavior {
	return &telemetry{}
}

func (app *telemetry) Load(_ gen.Node, _ ...any) (spec gen.ApplicationSpec, err error) {
	spec.Name = "telemetry_app"
	spec.Description = "TELEMETRY"
	spec.Mode = gen.ApplicationModeTransient
	spec.Group = []gen.ApplicationMemberSpec{
		{
			Name:    "telemetry_state",
			Factory: newState,
		},
	}

	return spec, err
}

func (app *telemetry) Start(_ gen.ApplicationMode) {}
func (app *telemetry) Terminate(_ error)           {}
